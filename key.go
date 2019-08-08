package kubelock

func lockAnnotation(name string) string {
	return "kubelock.giantswarm.io/" + name
}
