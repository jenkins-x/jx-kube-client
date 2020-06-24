package kubeclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// this is so that we load the auth plugins so we can connect to, say, GCP
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	// PodNamespaceFile the file path and name for pod namespace
	PodNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	DefaultKubeConfigFile = "config"
	DefaultKubeConfigPath = ".kube"
)

type factory struct {
}

// NewFactory creates a factory with the default Kubernetes resources defined
func NewFactory() Factory {
	f := &factory{}
	return f
}

func (f *factory) CreateKubeConfigFromCustomLocation(kubeConfigPath, kubeConfigFile string) (*rest.Config, error) {
	return f.createKubeConfig(kubeConfigPath, kubeConfigFile)
}

func (f *factory) CreateKubeConfig() (*rest.Config, error) {
	return f.createKubeConfig(DefaultKubeConfigPath, DefaultKubeConfigFile)
}

// CreateKubeConfig figures out the kubernetes config from environment variables or default locations whether in or out
// of cluster
func (f *factory) createKubeConfig(kubeConfigPath, kubeConfigFile string) (*rest.Config, error) {
	masterURL := ""
	kubeConfigEnv := os.Getenv("KUBECONFIG")
	if kubeConfigEnv != "" {
		pathList := filepath.SplitList(kubeConfigEnv)
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{Precedence: pathList},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterURL}}).ClientConfig()
	}
	kubeconfig := f.createKubeConfigPath(kubeConfigPath, kubeConfigFile)
	var config *rest.Config
	var err error
	if kubeconfig != nil {
		exists, err := fileExists(*kubeconfig)
		if err == nil && exists {
			// use the current context in kubeconfig
			config, err = clientcmd.BuildConfigFromFlags(masterURL, *kubeconfig)
			if err != nil {
				return nil, err
			}
		}
	}
	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	// for testing purposes one can enable tracing of Kube REST API calls
	traceKubeAPI := os.Getenv("TRACE_KUBE_API")
	if traceKubeAPI == "1" || traceKubeAPI == "on" {
		config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return &Tracer{RoundTripper: rt}
		}
	}
	return config, nil
}

func (f *factory) createKubeConfigPath(kubeConfigPath, kubeConfigFile string) *string {
	path := ""
	if home := homeDir(); home != "" {
		path = filepath.Join(home, kubeConfigPath, kubeConfigFile)
	}
	return &path
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, errors.Wrapf(err, "failed to check if file exists %s", path)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h := os.Getenv("USERPROFILE") // windows
	if h == "" {
		h = "."
	}
	return h
}

// CurrentContext returns the current context
func CurrentContext(config *api.Config) *api.Context {
	if config != nil {
		name := config.CurrentContext
		if name != "" && config.Contexts != nil {
			return config.Contexts[name]
		}
	}
	return nil
}

// CurrentNamespace returns the current namespace in the context
func CurrentNamespace() (string, error) {
	config, _, err := LoadConfig()
	if err != nil {
		return "", err
	}
	ctx := CurrentContext(config)
	if ctx != nil {
		n := ctx.Namespace
		if n != "" {
			return n, nil
		}
	}
	// if we are in a pod lets try load the pod namespace file
	data, err := ioutil.ReadFile(PodNamespaceFile)
	if err == nil {
		n := string(data)
		if n != "" {
			return n, nil
		}
	}
	return "default", nil
}

// LoadConfig loads the Kubernetes configuration
func LoadConfig() (*api.Config, *clientcmd.PathOptions, error) {
	po := clientcmd.NewDefaultPathOptions()
	if po == nil {
		return nil, po, fmt.Errorf("could not find any default path options for the kubeconfig file usually found at ~/.kube/config")
	}
	config, err := po.GetStartingConfig()
	if err != nil {
		return nil, po, fmt.Errorf("could not load the kube config file %s due to %s", po.GetDefaultFilename(), err)
	}
	return config, po, err
}
