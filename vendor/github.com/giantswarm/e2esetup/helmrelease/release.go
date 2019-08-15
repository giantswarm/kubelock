package helmrelease

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/giantswarm/apprclient"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/e2esetup/internal/filelogger"
	"github.com/giantswarm/e2esetup/k8s"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/helm"
)

const (
	defaultNamespace = "default"
)

type Config struct {
	ApprClient *apprclient.Client
	HelmClient *helmclient.Client
	Logger     micrologger.Logger
	K8sClients *k8s.Clients

	Namespace string
}

type Release struct {
	apprClient *apprclient.Client
	helmClient *helmclient.Client
	logger     micrologger.Logger
	k8sClients *k8s.Clients

	condition  *conditionSet
	fileLogger *filelogger.FileLogger

	namespace string
}

func New(config Config) (*Release, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.ApprClient == nil {
		config.Logger.Log("level", "debug", "message", fmt.Sprintf("%T.ApprClient is empty", config))

		config.Logger.Log("level", "debug", "message", fmt.Sprintf("using default for %T.ApprClient", config))

		c := apprclient.Config{
			Fs:     afero.NewOsFs(),
			Logger: config.Logger,

			Address:      "https://quay.io",
			Organization: "giantswarm",
		}

		a, err := apprclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		config.ApprClient = a
	}
	if config.HelmClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HelmClient must not be empty", config)
	}
	if config.K8sClients == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClients must not be empty", config)
	}
	if config.Namespace == "" {
		config.Namespace = defaultNamespace
	}

	var err error

	var condition *conditionSet
	{

		c := conditionSetConfig{
			K8sClients: config.K8sClients,
			Logger:     config.Logger,
		}

		condition, err = newConditionSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var fileLogger *filelogger.FileLogger
	{
		c := filelogger.Config{
			K8sClient: config.K8sClients.K8sClient(),
			Logger:    config.Logger,
		}

		fileLogger, err = filelogger.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	r := &Release{
		apprClient: config.ApprClient,
		helmClient: config.HelmClient,
		k8sClients: config.K8sClients,
		logger:     config.Logger,

		namespace: config.Namespace,

		condition:  condition,
		fileLogger: fileLogger,
	}

	return r, nil
}

func (r *Release) Condition() ConditionSet {
	return r.condition
}

func (r *Release) Delete(ctx context.Context, name string) error {
	releaseName := fmt.Sprintf("%s-%s", r.namespace, name)

	err := r.helmClient.DeleteRelease(ctx, releaseName, helm.DeletePurge(true))
	if helmclient.IsReleaseNotFound(err) {
		return microerror.Maskf(releaseNotFoundError, "failed to delete release %#q", name)
	} else if helmclient.IsTillerNotFound(err) {
		return microerror.Maskf(tillerNotFoundError, "failed to delete release %#q", name)
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// EnsureDeleted makes sure the release is deleted and purged and all
// conditions are met.
func (r *Release) EnsureDeleted(ctx context.Context, name string, conditions ...ConditionFunc) error {
	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting release %#q", name))

		err := r.Delete(ctx, name)
		if IsReleaseNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("release %#q already deleted", name))
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("deleted release %#q", name))
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured deletion of release %#q", name))
	}

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensuring conditions for deleted release %#q", name))

		err := r.waitForConditions(ctx, conditions...)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured conditions for deleted release %#q", name))
	}

	return nil
}

// EnsureInstalled makes sure the release is installed and all conditions are
// met. If release name ends with "-operator" suffix it also selects
// a "app=${name}" pod and streams it logs to the ./logs directory.
//
// NOTE: It does not update the release if it already exists.
func (r *Release) EnsureInstalled(ctx context.Context, name string, chartInfo ChartInfo, values string, conditions ...ConditionFunc) error {
	var err error
	isOperator := strings.HasSuffix(name, "-operator")

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating release %#q", name))

		err := r.Install(ctx, name, chartInfo, values)
		if IsReleaseAlreadyExists(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("release %#q already created", name))
		} else if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created release %#q", name))
	}

	if isOperator {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensuring operator pod exists for release %#q", name))

		c := r.Condition().PodExists(ctx, r.namespace, fmt.Sprintf("app=%s", name))
		err := r.waitForConditions(ctx, c)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured operator pod exists for release %#q", name))
	}

	var operatorPodName string
	var operatorPodNamespace string = r.namespace
	if isOperator {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("finding operator pod name for release %#q", name))

		operatorPodName, err = r.podName(operatorPodNamespace, fmt.Sprintf("app=%s", name))
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("found operator pod name for release %#q", name))
	}

	if isOperator {
		err := r.fileLogger.EnsurePodLogging(ctx, operatorPodNamespace, operatorPodName)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensuring conditions for release %#q", name))

		err := r.waitForConditions(ctx, conditions...)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured conditions for release %#q", name))
	}

	return nil
}

func (r *Release) Install(ctx context.Context, name string, chartInfo ChartInfo, values string, conditions ...ConditionFunc) error {
	releaseName := fmt.Sprintf("%s-%s", r.namespace, name)

	var err error

	tarballPath, err := r.pullTarball(ctx, name, chartInfo)
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.helmClient.InstallReleaseFromTarball(ctx, tarballPath, r.namespace, helm.ReleaseName(releaseName), helm.ValueOverrides([]byte(values)), helm.InstallWait(true))
	if helmclient.IsReleaseAlreadyExists(err) {
		return microerror.Maskf(releaseAlreadyExistsError, "failed to install release %#q", releaseName)
	} else if helmclient.IsTarballNotFound(err) {
		return microerror.Maskf(tarballNotFoundError, "failed to install release %#q from tarball %#q", releaseName, tarballPath)
	} else if err != nil {
		return microerror.Mask(err)
	}

	err = r.waitForConditions(ctx, conditions...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Release) Update(ctx context.Context, name string, chartInfo ChartInfo, values string, conditions ...ConditionFunc) error {
	releaseName := fmt.Sprintf("%s-%s", r.namespace, name)

	var err error

	tarballPath, err := r.pullTarball(ctx, name, chartInfo)
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.helmClient.UpdateReleaseFromTarball(ctx, releaseName, tarballPath, helm.UpdateValueOverrides([]byte(values)), helm.UpgradeWait(true))
	if helmclient.IsReleaseAlreadyExists(err) {
		return microerror.Maskf(releaseAlreadyExistsError, "failed to update release %#q", releaseName)
	} else if helmclient.IsTarballNotFound(err) {
		return microerror.Maskf(tarballNotFoundError, "failed to update release %#q from tarball %#q", releaseName, tarballPath)
	} else if err != nil {
		return microerror.Mask(err)
	}

	err = r.waitForConditions(ctx, conditions...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Release) WaitForStatus(ctx context.Context, release string, status string) error {
	operation := func() error {
		rc, err := r.helmClient.GetReleaseContent(ctx, release)
		if helmclient.IsReleaseNotFound(err) && status == "DELETED" {
			// Error is expected because we purge releases when deleting.
			return nil
		} else if err != nil {
			return microerror.Mask(err)
		}
		if rc.Status != status {
			return microerror.Maskf(releaseStatusNotMatchingError, "waiting for '%s', current '%s'", status, rc.Status)
		}
		return nil
	}

	notify := func(err error, t time.Duration) {
		r.logger.Log("level", "debug", "message", fmt.Sprintf("failed to get release status '%s': retrying in %s", status, t), "stack", fmt.Sprintf("%v", err))
	}

	b := backoff.NewExponential(backoff.MediumMaxWait, backoff.LongMaxInterval)
	err := backoff.RetryNotify(operation, b, notify)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (r *Release) WaitForChartInfo(ctx context.Context, release string, version string) error {
	operation := func() error {
		rh, err := r.helmClient.GetReleaseHistory(ctx, release)
		if err != nil {
			return microerror.Mask(err)
		}
		if rh.Version != version {
			return microerror.Maskf(releaseVersionNotMatchingError, "waiting for '%s', current '%s'", version, rh.Version)
		}
		return nil
	}

	notify := func(err error, t time.Duration) {
		r.logger.Log("level", "debug", "message", fmt.Sprintf("failed to get release version '%s': retrying in %s", version, t), "stack", fmt.Sprintf("%v", err))
	}

	b := backoff.NewExponential(backoff.ShortMaxWait, backoff.LongMaxInterval)
	err := backoff.RetryNotify(operation, b, notify)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (r *Release) podName(namespace, labelSelector string) (string, error) {
	pods, err := r.k8sClients.K8sClient().CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return "", microerror.Mask(err)
	}
	if len(pods.Items) > 1 {
		return "", microerror.Maskf(waitError, "expected at most 1 pod but got %d", len(pods.Items))
	}
	if len(pods.Items) == 0 {
		return "", microerror.Mask(notFoundError)
	}
	pod := pods.Items[0]
	return pod.Name, nil
}

func (r *Release) pullTarball(ctx context.Context, releaseName string, chartInfo ChartInfo) (string, error) {
	chartName := chartInfo.name
	if chartName == "" {
		chartName = fmt.Sprintf("%s-chart", releaseName)
	}

	if chartInfo.isChannel {
		tarball, err := r.apprClient.PullChartTarball(ctx, chartName, chartInfo.version)
		if err != nil {
			return "", microerror.Mask(err)
		}

		return tarball, nil
	}

	tarball, err := r.apprClient.PullChartTarballFromRelease(ctx, chartName, chartInfo.version)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return tarball, nil
}

func (r *Release) waitForConditions(ctx context.Context, conditions ...ConditionFunc) error {
	for _, c := range conditions {
		err := ctx.Err()
		if err != nil {
			return microerror.Mask(err)
		}

		o := func() error {
			err := c()
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.MediumMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(r.logger, ctx)

		err = backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
