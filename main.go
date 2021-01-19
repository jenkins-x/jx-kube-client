package main

import (
	"context"
	"fmt"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

func main() {
	kcfg := os.Getenv("KUBECONFIG")
	fmt.Printf("config: %s\n", kcfg)

	f := kubeclient.NewFactory()
	cfg, err := f.CreateKubeConfig()
	if err != nil {
		fmt.Printf("failed to create config: %s\n", err.Error())
		return
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Printf("failed to create client: %s\n", err.Error())
		return
	}

	ctx := context.TODO()
	ns, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("failed to list namespaces: %s\n", err.Error())
		return
	}

	fmt.Printf("has %d namespaces:\n", len(ns.Items))

	for _, n := range ns.Items {
		fmt.Printf("  %s\n", n.Name)
	}
}
