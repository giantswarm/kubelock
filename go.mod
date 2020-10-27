module github.com/giantswarm/kubelock/v3

go 1.14

require (
	github.com/giantswarm/k8sclient/v5 v5.0.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.3
	github.com/google/go-cmp v0.5.2
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
)

replace sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
