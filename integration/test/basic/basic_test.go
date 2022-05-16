//go:build k8srequired
// +build k8srequired

package basic

import (
	"context"
	"testing"

	"github.com/giantswarm/kubelock/v4"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestBasic_ClusterScope(t *testing.T) {
	var err error
	ctx := context.Background()

	var kubeLock *kubelock.KubeLock
	{
		c := kubelock.Config{
			DynClient: config.K8sClients.DynClient(),

			GVR: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "namespaces",
			},
		}

		kubeLock, err = kubelock.New(c)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	lock := kubeLock.Lock("test-lock")

	err = lock.Acquire(ctx, "default", kubelock.AcquireOptions{})
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.JSON(err))
	}

	err = lock.Acquire(ctx, "default", kubelock.AcquireOptions{})
	if !kubelock.IsAlreadyExists(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
	}

	err = lock.Release(ctx, "default", kubelock.ReleaseOptions{})
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.JSON(err))
	}

	err = lock.Release(ctx, "default", kubelock.ReleaseOptions{})
	if !kubelock.IsNotFound(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.JSON(err))
	}
}

func TestBasic_Namespaced(t *testing.T) {
	var err error
	ctx := context.Background()

	var kubeLock *kubelock.KubeLock
	{
		c := kubelock.Config{
			DynClient: config.K8sClients.DynClient(),

			GVR: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "configmaps",
			},
		}

		kubeLock, err = kubelock.New(c)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	lock := kubeLock.Lock("test-lock")

	err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", kubelock.AcquireOptions{})
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.JSON(err))
	}

	err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", kubelock.AcquireOptions{})
	if !kubelock.IsAlreadyExists(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
	}

	err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", kubelock.ReleaseOptions{})
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.JSON(err))
	}

	err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", kubelock.ReleaseOptions{})
	if !kubelock.IsNotFound(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.JSON(err))
	}
}
