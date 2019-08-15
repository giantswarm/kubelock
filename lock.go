package kubelock

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

type lock struct {
	resource dynamic.ResourceInterface

	lockName string
}

func (l *lock) Acquire(ctx context.Context, name string) error {
	obj, err := l.resource.Get(name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Check if there is lock acquired and error if so.
	{
		ann := obj.GetAnnotations()
		_, ok := ann[lockAnnotation(l.lockName)]
		if ok {
			return microerror.Maskf(alreadyExistsError, "lock %#q on %#q already acquired", l.lockName, obj.GetSelfLink())
		}
	}

	var data []byte
	{
		// To simplify the PR I removed TTL and owner handling. That
		// will come in the follow up PR.
		data = []byte("TODO")
	}

	// Update object annotations.
	{
		ann := obj.GetAnnotations()
		if ann == nil {
			ann = map[string]string{}
		}
		ann[lockAnnotation(l.lockName)] = string(data)
		obj.SetAnnotations(ann)
	}

	// Update object.
	{
		_, err := l.resource.Update(obj, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (l *lock) Release(ctx context.Context, name string) error {
	obj, err := l.resource.Get(name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Check if the lock exists and fail if it doesn't.
	{
		ann := obj.GetAnnotations()
		_, ok := ann[lockAnnotation(l.lockName)]
		if !ok {
			return microerror.Maskf(notFoundError, "lock %#q on %#q not found", l.lockName, obj.GetSelfLink())
		}
	}

	// Update object annotations.
	{
		ann := obj.GetAnnotations()
		delete(ann, lockAnnotation(l.lockName))
		obj.SetAnnotations(ann)
	}

	// Update object.
	{
		_, err := l.resource.Update(obj, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (l *lock) data(obj *unstructured.Unstructured) (string, bool, error) {
	ann := obj.GetAnnotations()
	stringData, ok := ann[lockAnnotation(l.lockName)]
	if !ok {
		return "", false, nil
	}

	// The name stringData and returning error is weird but it will make
	// more sense when there is actual data.
	return stringData, true, nil
}
