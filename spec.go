package kubelock

import (
	"context"
	"time"
)

const (
	// DefaultTTL is default time to live for the lock.
	DefaultTTL = 5 * time.Minute
)

// Interface is the interface of a distributed Kubernetes lock. The default
// implementation is KubeLock.
//
// The typical usage for a namespace resource may look like:
//
//	kubeLock.Lock("my-lock-name").Namespace("my-namespace").Acquire(ctx, "my-configmap", kubelock.LockOptions{})
//
// The typical usage for a cluster scope resource may look like:
//
//	kubeLock.Lock("my-lock-name").Acquire(ctx, "my-namespace", kubelock.LockOptions{})
//
type Interface interface {
	// Lock creates a lock with the given name. The name will be used to
	// create annotation prefixed with "kubelock.giantswarm.io/" on the
	// Kubernetes object. Value of this annotation stores the lock data.
	Lock(name string) NamespaceableLock
}

type Lock interface {
	// Acquire tries to acquire the lock on a Kubernetes resource with the
	// given name.
	Acquire(ctx context.Context, name string, options LockOptions) error
	// Release tries to release the lock on a Kubernetes resource with the
	// given name.
	Release(ctx context.Context, name string, options LockOptions) error
}

type NamespaceableLock interface {
	Lock

	// Namespace creates a lock that can be acquired on Kubernetes
	// resources in the given namespace.
	Namespace(ns string) Lock
}
