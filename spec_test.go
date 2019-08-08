package kubelock

import "testing"

func Test_Interface(t *testing.T) {
	var _ Interface = &KubeLock{}
}
