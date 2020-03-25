package kubelock

import "github.com/giantswarm/microerror"

var alreadyExistsError = &microerror.Error{
	Kind: "alreadyExistsError",
}

// IsAlreadyExists asserts alreadyExistsError.
func IsAlreadyExists(err error) bool {
	return microerror.Cause(err) == alreadyExistsError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var ownerMismatchError = &microerror.Error{
	Kind: "ownerMismatchError",
}

// IsOwnerMismatch asserts ownerMismatchError.
func IsOwnerMismatch(err error) bool {
	return microerror.Cause(err) == ownerMismatchError
}
