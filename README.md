# jx-kube-client

Used to create a kubernetes client from within a Kubernetes cluster or ouside using ~/.kube/config

Here's an example which also uses [Jenkins X logging](https://github.com/jenkins-x/jx-logging)

```go
import (
    "github.com/jenkins-x/jx-kube-client/pkg/kubeclient"
    "github.com/jenkins-x/jx-logging/pkg/log"
)

func main() {
    f := kubeclient.NewFactory()
    cfg, err := f.CreateKubeConfig()
    if err != nil {
        log.Logger().Fatalf("failed to get kubernetes config: %v", err)
    }


    kubeClient, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        log.Logger().Fatalf("error building kubernetes clientset: %v", err)
    }
}
```


Part of Jenkins X shared libraries.