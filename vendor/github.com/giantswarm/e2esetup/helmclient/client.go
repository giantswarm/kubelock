package helmclient

import (
	helmclientlib "github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClientConfig struct {
	K8sClient  kubernetes.Interface
	Logger     micrologger.Logger
	RestConfig *rest.Config

	TillerImage     string
	TillerNamespace string
}

type Client struct {
	logger micrologger.Logger

	helmClient helmclientlib.Interface
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.RestConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RestConfig must not be empty", config)
	}

	var err error

	var helmClient helmclientlib.Interface
	{
		c := helmclientlib.Config{
			K8sClient:  config.K8sClient,
			Logger:     config.Logger,
			RestConfig: config.RestConfig,

			TillerImage:     config.TillerImage,
			TillerNamespace: config.TillerNamespace,
		}

		helmClient, err = helmclientlib.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Client{
		logger: config.Logger,

		helmClient: helmClient,
	}

	return c, nil
}

func (c *Client) HelmClient() helmclientlib.Interface {
	return c.helmClient
}
