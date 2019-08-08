package kubelock

import (
	"time"
)

func lockAnnotation(name string) string {
	return "kubelock.giantswarm.io/" + name
}

func isExpired(data lockData) bool {
	return data.CreatedAt.Add(data.TTL).Before(time.Now())
}

func lockAnnotation(name string) string {
	return "kubelock.giantswarm.io/" + name
}

func defaultedOptions(options LockOptions) LockOptions {
	if options.TTL == 0 {
		options.TTL = DefaultTTL
	}

	return options
}
