package helmrelease

import (
	"context"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/e2esetup/k8s"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	v1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type conditionSetConfig struct {
	K8sClients *k8s.Clients
	Logger     micrologger.Logger
}

type conditionSet struct {
	k8sClients *k8s.Clients
	logger     micrologger.Logger
}

func newConditionSet(config conditionSetConfig) (*conditionSet, error) {
	if config.K8sClients == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClients must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	c := &conditionSet{
		k8sClients: config.K8sClients,
		logger:     config.Logger,
	}

	return c, nil
}

func (c *conditionSet) CRDExists(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition) ConditionFunc {
	return func() error {
		o := func() error {
			_, err := c.k8sClients.ExtClient().ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return microerror.Mask(err)
			} else if err != nil {
				return backoff.Permanent(microerror.Mask(err))
			}
			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(c.logger, ctx)
		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (c *conditionSet) CRDNotFound(ctx context.Context, crd *apiextensionsv1beta1.CustomResourceDefinition) ConditionFunc {
	return func() error {
		o := func() error {
			_, err := c.k8sClients.ExtClient().ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return backoff.Permanent(microerror.Mask(err))
			}

			return microerror.Maskf(waitError, "CRD %#q still exists", crd.Name)
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(c.logger, ctx)
		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (c *conditionSet) PodExists(ctx context.Context, namespace, labelSelector string) ConditionFunc {
	return func() error {
		o := func() error {
			pods, err := c.k8sClients.K8sClient().CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return microerror.Mask(err)
			}
			if len(pods.Items) != 1 {
				return microerror.Maskf(waitError, "expected 1 pod but got %d", len(pods.Items))
			}

			pod := pods.Items[0]
			if pod.Status.Phase != v1.PodRunning {
				return microerror.Maskf(waitError, "expected Pod phase %#q but got %#q", v1.PodRunning, pod.Status.Phase)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.MediumMaxWait, backoff.LongMaxInterval)

		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (c *conditionSet) PodNotFound(ctx context.Context, namespace, labelSelector string) ConditionFunc {
	return func() error {
		o := func() error {
			pods, err := c.k8sClients.K8sClient().CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
			if err != nil {
				return microerror.Mask(err)
			}

			if len(pods.Items) != 0 {
				return microerror.Maskf(waitError, "expected no Pods for label selector %#q but got %d", labelSelector, len(pods.Items))
			}

			return nil
		}
		b := backoff.NewExponential(backoff.MediumMaxWait, backoff.LongMaxInterval)

		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (c *conditionSet) SecretExists(ctx context.Context, namespace, name string) ConditionFunc {
	return func() error {
		o := func() error {
			_, err := c.k8sClients.K8sClient().CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return microerror.Mask(err)
			} else if err != nil {
				return backoff.Permanent(microerror.Mask(err))
			}
			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(c.logger, ctx)
		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}

func (c *conditionSet) SecretNotFound(ctx context.Context, namespace, name string) ConditionFunc {
	return func() error {
		o := func() error {
			_, err := c.k8sClients.K8sClient().CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return backoff.Permanent(microerror.Mask(err))
			}
			return microerror.Maskf(waitError, "Secret %#q in namespace %#q still exists", name, namespace)
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(c.logger, ctx)

		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}
}
