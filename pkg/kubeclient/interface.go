package kubeclient

import (
	"k8s.io/client-go/rest"
)

// Factory is the interface defined for Kubernetes, Jenkins X, and Tekton REST APIs
type Factory interface {
	// CreateKubeConfig creates the kubernetes configuration
	CreateKubeConfig() (*rest.Config, error)
	CreateKubeConfigFromCustomLocation(string, string) (*rest.Config, error)
}
