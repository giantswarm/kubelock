package kubelock

import (
	"time"
)

type LockOptions struct {
	TTL time.Duration
}

type lockData struct {
	CreatedAt time.Time     `json:"createdAt"`
	TTL       time.Duration `json:"ttl"`
}

type jsonPatch []jsonPatchOperation

type jsonPatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func newAcquirePatch(resourceVersion, lockName string, lockData []byte) jsonPatch {
	return jsonPatch{
		{
			Op:    "test",
			Path:  "/metadata/resourceVersion",
			Value: resourceVersion,
		},
		{
			Op:    "add",
			Path:  "/metadata/annotations/" + lockAnnotation(lockName),
			Value: string(lockData),
		},
	}
}

func newReleasePatch(resourceVersion, lockName string) jsonPatch {
	return jsonPatch{
		{
			Op:    "test",
			Path:  "/metadata/resourceVersion",
			Value: resourceVersion,
		},
		{
			Op:   "remove",
			Path: "/metadata/annotations/" + lockAnnotation(lockName),
		},
	}
}
