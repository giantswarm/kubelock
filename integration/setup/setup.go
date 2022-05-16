//go:build k8srequired
// +build k8srequired

package setup

import (
	"context"
	"os"
	"testing"

	"github.com/giantswarm/microerror"
)

func Setup(m *testing.M, config Config) {
	ctx := context.Background()

	v, err := setup(ctx, m, config)
	if err != nil {
		config.Logger.LogCtx(ctx, "level", "error", "message", "failed to setup test environment", "stack", microerror.JSON(err))
		os.Exit(1)
	}

	os.Exit(v)
}

func setup(ctx context.Context, m *testing.M, config Config) (int, error) {
	// Create namespace.
	{
		err := config.K8sSetup.EnsureNamespaceCreated(ctx, "testing")
		if err != nil {
			return 0, microerror.Mask(err)
		}
	}

	v := m.Run()
	return v, nil
}
