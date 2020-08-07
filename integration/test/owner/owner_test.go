// +build k8srequired

package owner

import (
	"context"
	"testing"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/giantswarm/kubelock/v2"
)

func TestOwner_ClusterScope(t *testing.T) {
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
	owner := "test-owner"

	// Acquire the lock.
	{
		opts := kubelock.AcquireOptions{
			Owner: owner,
		}
		err = lock.Acquire(ctx, "default", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Try to acquire the lock again with the same owner and check for
	// already exists error.
	{
		opts := kubelock.AcquireOptions{
			Owner: owner,
		}
		err = lock.Acquire(ctx, "default", opts)
		if !kubelock.IsAlreadyExists(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
		}
	}

	// Try to acquire the lock again with a different owner and check for
	// owner mismatch error.
	{
		opts := kubelock.AcquireOptions{
			Owner: owner + "-i-am-a-different-owner",
		}
		err = lock.Acquire(ctx, "default", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
		}
	}

	// Try to acquire the lock again with a no owner and check for already
	// owner mismatch error.
	{
		opts := kubelock.AcquireOptions{
			Owner: "",
		}
		err = lock.Acquire(ctx, "default", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
		}
	}

	// Try to release the lock with a different owner and check if it fails
	// with owner mismatch error.
	{
		opts := kubelock.ReleaseOptions{
			Owner: owner + "-i-am-a-different-owner",
		}
		err = lock.Release(ctx, "default", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsOwnerMismatch", microerror.JSON(err))
		}
	}

	// Try to release the lock with no owner and check if it fails with
	// owner mismatch error.
	{
		opts := kubelock.ReleaseOptions{
			Owner: "",
		}
		err = lock.Release(ctx, "default", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsOwnerMismatch", microerror.JSON(err))
		}
	}

	// Release the lock with the same owner.
	{
		opts := kubelock.ReleaseOptions{
			Owner: owner,
		}
		err = lock.Release(ctx, "default", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}
}

func TestOwner_Namespaced(t *testing.T) {
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
	owner := "test-owner"

	// Acquire the lock.
	{
		opts := kubelock.AcquireOptions{
			Owner: owner,
		}
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}

	// Try to acquire the lock again with the same owner and check for
	// already exists error.
	{
		opts := kubelock.AcquireOptions{
			Owner: owner,
		}
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", opts)
		if !kubelock.IsAlreadyExists(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
		}
	}

	// Try to acquire the lock again with a different owner and check for
	// owner mismatch error.
	{
		opts := kubelock.AcquireOptions{
			Owner: owner + "-i-am-a-different-owner",
		}
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
		}
	}

	// Try to acquire the lock again with a no owner and check for already
	// owner mismatch error.
	{
		opts := kubelock.AcquireOptions{
			Owner: "",
		}
		err = lock.Namespace("kube-system").Acquire(ctx, "kube-proxy", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsAlreadyExists", microerror.JSON(err))
		}
	}

	// Try to release the lock with a different owner and check if it fails
	// with owner mismatch error.
	{
		opts := kubelock.ReleaseOptions{
			Owner: owner + "-i-am-a-different-owner",
		}
		err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsOwnerMismatch", microerror.JSON(err))
		}
	}

	// Try to release the lock with no owner and check if it fails with
	// owner mismatch error.
	{
		opts := kubelock.ReleaseOptions{
			Owner: "",
		}
		err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", opts)
		if !kubelock.IsOwnerMismatch(err) {
			t.Fatalf("error == %#v, want matching kubelock.IsOwnerMismatch", microerror.JSON(err))
		}
	}

	// Release the lock with the same owner.
	{
		opts := kubelock.ReleaseOptions{
			Owner: owner,
		}
		err = lock.Namespace("kube-system").Release(ctx, "kube-proxy", opts)
		if err != nil {
			t.Fatalf("error == %#q, want nil", microerror.JSON(err))
		}
	}
}
