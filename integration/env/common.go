package env

import (
	"os"
)

const (
	EnvVarE2EKubeconfig = "E2E_KUBECONFIG"
)

var (
	kubeconfig string
)

func init() {
	kubeconfig = os.Getenv(EnvVarE2EKubeconfig)
}

func KubeConfigPath() string {
	return kubeconfig
}
