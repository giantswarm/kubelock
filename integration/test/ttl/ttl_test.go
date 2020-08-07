// +build k8srequired

package ttl

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/kubelock/v2"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestTTL_ClusterScope_Acquire(t *testing.T) {
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

	lock := kubeLock.Lock("test-lock-acquire")
	ttl := 100 * time.Millisecond

	// Acquire the lock.
	{
		opts := kubelock.AcquireOptions{
			TTL: ttl,
		}
		err = lock.Acquire(ctx, "default", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Wait for the lock to expire
	{
		time.Sleep(ttl + 10*time.Millisecond)
	}

	// Check if the lock can be acquired again.
	{
		err = lock.Acquire(ctx, "default", kubelock.AcquireOptions{})
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Clean up.
	{
		err = lock.Release(ctx, "default", kubelock.ReleaseOptions{})
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}
}

func TestTTL_ClusterScope_Release(t *testing.T) {
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

	lock := kubeLock.Lock("test-lock-release")
	ttl := 100 * time.Millisecond

	// Acquire the lock.
	{
		opts := kubelock.AcquireOptions{
			TTL: ttl,
		}
		err = lock.Acquire(ctx, "default", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Wait for the lock to expire
	{
		time.Sleep(ttl + 10*time.Millisecond)
	}

	// Check if the lock is not found.
	{
		err = lock.Release(ctx, "default", kubelock.ReleaseOptions{})
		if !kubelock.IsNotFound(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.JSON(err))
		}
	}
}

func TestBasic_Namespaced_Acquire(t *testing.T) {
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

	lock := kubeLock.Lock("test-lock-acquire")
	ttl := 100 * time.Millisecond

	// Acquire the lock.
	{
		opts := kubelock.AcquireOptions{
			TTL: ttl,
		}
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Wait for the lock to expire
	{
		time.Sleep(ttl + 10*time.Millisecond)
	}

	// Check if the lock is not found.
	{
		err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", kubelock.ReleaseOptions{})
		if !kubelock.IsNotFound(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.JSON(err))
		}
	}

	// Check if the lock can be acquired again.
	{
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", kubelock.AcquireOptions{})
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Clean up.
	{
		err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", kubelock.ReleaseOptions{})
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}
}

func TestBasic_Namespaced_Release(t *testing.T) {
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

	lock := kubeLock.Lock("test-lock-release")
	ttl := 100 * time.Millisecond

	// Acquire the lock.
	{
		opts := kubelock.AcquireOptions{
			TTL: ttl,
		}
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Wait for the lock to expire
	{
		time.Sleep(ttl + 10*time.Millisecond)
	}

	// Check if the lock is not found.
	{
		err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", kubelock.ReleaseOptions{})
		if !kubelock.IsNotFound(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsNotFound", microerror.JSON(err))
		}
	}
}
