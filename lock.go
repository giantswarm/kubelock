package kubelock

import (
	"context"
	"encoding/json"
	"time"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

type lock struct {
	resource dynamic.ResourceInterface

	lockName string
}

func (l *lock) Acquire(ctx context.Context, name string, options LockOptions) error {
	options = defaultedOptions(options)

	obj, err := l.resource.Get(name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Check if there is non expired lock acquired and error if so.
	{
		data, ok, err := l.data(obj)
		if err != nil {
			return microerror.Mask(err)
		}
		if ok && !isExpired(data) {
			if data.Owner == options.Owner {
				return microerror.Maskf(alreadyExistsError, "lock %#q on %#q owned by %#q already acquired at %s with TTL %s", l.lockName, obj.GetSelfLink(), data.Owner, data.CreatedAt.Format(time.RFC3339), data.TTL)
			} else {
				return microerror.Maskf(ownerMismatchError, "lock %#q on %#q owned by %#q already acquired at %s with TTL %s", l.lockName, obj.GetSelfLink(), data.Owner, data.CreatedAt.Format(time.RFC3339), data.TTL)
			}
		}
	}

	var data []byte
	{
		d := lockData{
			Onwer:     options.Owner,
			CreatedAt: time.Now(),
			TTL:       options.TTL,
		}

		data, err = json.Marshal(d)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var patch []byte
	{
		p := newAcquirePatch(obj.GetResourceVersion(), l.lockName, data)

		patch, err = json.Marshal(p)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	_, err = l.resource.Patch(name, types.JSONPatchType, patch, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (l *lock) Release(ctx context.Context, name string, options LockOptions) error {
	options = defaultedOptions(options)

	obj, err := l.resource.Get(name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Check if the lock exists and fail if it doesn't.
	{
		data, ok, err := l.data(obj)
		if err != nil {
			return microerror.Mask(err)
		}
		if !ok || isExpired(data) {
			return microerror.Maskf(notFoundError, "lock %#q on %#q not found", l.lockName, obj.GetSelfLink())
		}
		if ok && !isExpired(data) && data.Owner != options.Owner {
			return microerror.Maskf(ownerMismatchError, "lock %#q on %#q is not owned by %#q", l.lockName, obj.GetSelfLink(), options.Owner)
		}
	}

	var patch []byte
	{
		p := newReleasePatch(obj.GetResourceVersion(), l.lockName)

		bs, err := json.Marshal(p)
		if err != nil {
			return microerror.Mask(err)
		}

		patch = bs
	}

	_, err = l.resource.Patch(name, types.JSONPatchType, patch, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (l *lock) data(obj *unstructured.Unstructured) (lockData, bool, error) {
	ann := obj.GetAnnotations()
	stringData, ok := ann[lockAnnotation(l.lockName)]
	if !ok {
		return lockData{}, false, nil
	}

	var data lockData
	err := json.Unmarshal([]byte(stringData), &data)
	if err != nil {
		return lockData{}, false, microerror.Mask(err)
	}

	return data, true, nil
}
