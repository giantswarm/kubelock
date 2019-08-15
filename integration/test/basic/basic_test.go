// +build k8srequired

package basic

import (
	"context"
	"testing"

	"github.com/giantswarm/kubelock"

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
			t.Fatalf("error == %#q, want nil", microerror.Stack(err))
		}
	}

	lock := kubeLock.Lock("test-lock")

	err = lock.Acquire(ctx, "default")
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.Stack(err))
	}

	err = lock.Acquire(ctx, "default")
	if !kubelock.IsAlreadyExists(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.Stack(err))
	}

	err = lock.Release(ctx, "default")
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.Stack(err))
	}

	err = lock.Release(ctx, "default")
	if !kubelock.IsNotFound(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.Stack(err))
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
			t.Fatalf("error == %#q, want nil", microerror.Stack(err))
		}
	}

	lock := kubeLock.Lock("test-lock")

	err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy")
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.Stack(err))
	}

	err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy")
	if !kubelock.IsAlreadyExists(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.Stack(err))
	}

	err = lock.Namespace("kube-system").Release(ctx, "kube-proxy")
	if err != nil {
		t.Fatalf("error == %#q, want nil", microerror.Stack(err))
	}

	err = lock.Namespace("kube-system").Release(ctx, "kube-proxy")
	if !kubelock.IsNotFound(err) {
		t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.Stack(err))
	}
}
