package kubelock

import (
	"context"
	"time"
)

// TODO: Make the repo public.

// TODO: Write acknowledgments to github.com/pulcy/kube-lock/ in the readme.

// TODO: Use dynamic client instead.

type Interface interface {
	Acquire(ctx context.Context, funcs ObjectFuncs, ttl time.Duration) error
	Release(ctx context.Context, funcs ObjectFuncs, ttl time.Duration) error
}
