package kubelock

import "context"

type ObjectFuncs interface {
	Get(ctx context.Context) (ObjectMeta, error)
	Update(ctx context.Context, meta ObjectMeta) error
}
