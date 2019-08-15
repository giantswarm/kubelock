package kubelock

import (
	"time"
)

type AcquireOptions struct {
	TTL time.Duration
}

type ReleaseOptions struct {
}

type lockData struct {
	CreatedAt time.Time     `json:"createdAt"`
	TTL       time.Duration `json:"ttl"`
}
