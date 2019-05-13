package kubelock

type ObjectMeta struct {
	Annotations     map[string]string
	ResourceVersion string
}
