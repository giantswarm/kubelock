package env

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvVarCircleCI      = "CIRCLECI"
	EnvVarCircleSHA     = "CIRCLE_SHA1"
	EnvVarE2EKubeconfig = "E2E_KUBECONFIG"
	EnvVarE2ETestDir    = "E2E_TEST_DIR"
)

var (
	circleCI   string
	circleSHA  string
	kubeconfig string
	testDir    string
)

func init() {
	circleCI = os.Getenv(EnvVarCircleCI)
	kubeconfig = os.Getenv(EnvVarE2EKubeconfig)

	circleSHA = os.Getenv(EnvVarCircleSHA)
	if circleSHA == "" {
		panic(fmt.Sprintf("env var '%s' must not be empty", EnvVarCircleSHA))
	}

	testDir = os.Getenv(EnvVarE2ETestDir)
	if testDir == "" {
		panic(fmt.Sprintf("env var '%s' must not be empty", EnvVarE2ETestDir))
	}

}

func CircleCI() bool {
	return strings.ToLower(circleCI) == "true"
}

func CircleSHA() string {
	return circleSHA
}

func KubeConfigPath() string {
	return kubeconfig
}

func TestDir() string {
	return testDir
}
