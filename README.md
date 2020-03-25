[![GoDoc](https://godoc.org/github.com/giantswarm/kubelock?status.svg)](http://godoc.org/github.com/giantswarm/kubelock) [![CircleCI](https://circleci.com/gh/giantswarm/kubelock.svg?style=shield)](https://circleci.com/gh/giantswarm/kubelock)

# kubelock

Package kubelock provides functionality to create distributed locks on
arbitrary kubernetes resources. It is heavily inspired by [pulcy/kube-lock] but
uses [client-go] library and its dynamic client.

## Usage

At Giant Swarm we run multiple instances of the same operators in different
versions. Some actions performed by operators are not concurrent. Good example
is IP range allocation for a newly created cluster. Each lock created by
kubelock has a name and an owner. A custom name allows to create multiple locks
on the same Kubernetes resource. The owner string usually contains the version
of the operator acquiring the lock. That way the operator can know if the lock
was acquired by itself or other operator version.



## Integration Tests

You can simply create a [`kind`](https://github.com/kubernetes-sigs/kind/)
cluster to run the integration tests.

```
kind create cluster
```

The tests need to figure out how to connect to the Kubernetes cluster. Therefore
we need to set an environment variable pointing to your local kube config.

```
export E2E_KUBECONFIG=~/.kube/config
```

Now you can easily run the integration tests. Note that `-count=1` is the
idomatic way to not cache tests, which we do not want for integration tests.

```
go test -tags=k8srequired ./integration/test/<test-name> -count=1
```

Once you did your testing you may want to delete your local test cluster again.

```
kind delete cluster
```



[client-go]: https://github.com/kubernetes/client-go
[pulcy/kube-lock]: https://github.com/pulcy/kube-lock
