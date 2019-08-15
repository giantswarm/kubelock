// +build k8srequired

package basic

import (
	"fmt"
	"testing"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/kubelock/integration/setup"
)

var (
	config setup.Config
)

func init() {
	err := initMainTest()
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func initMainTest() error {
	var err error

	{
		config, err = setup.NewConfig()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	setup.Setup(m, config)
}
