// +build k8srequired

package setup

import (
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/e2esetup/k8s"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/kubelock/integration/env"
)

const (
	namespace = "giantswarm"
)

type Config struct {
	K8sClients *k8sclient.Clients
	K8sSetup   *k8s.Setup
	Logger     micrologger.Logger
}

func NewConfig() (Config, error) {
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var k8sClients *k8sclient.Clients
	{
		kubeConfigPath := env.KubeConfigPath()
		if kubeConfigPath == "" {
			kubeConfigPath = harness.DefaultKubeConfig
		}

		c := k8sclient.ClientsConfig{
			Logger:         logger,
			KubeConfigPath: kubeConfigPath,
		}

		k8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var k8sSetup *k8s.Setup
	{
		c := k8s.SetupConfig{
			Clients: k8sClients,
			Logger:  logger,
		}

		k8sSetup, err = k8s.NewSetup(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		K8sClients: k8sClients,
		K8sSetup:   k8sSetup,
		Logger:     logger,
	}

	return c, nil
}
