package server

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashfunc/debotops/api/v1alpha1"
	"github.com/hashfunc/debotops/pkg/core"
)

func getClientConfigLoader() clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := clientcmd.ConfigOverrides{}
	return clientcmd.NewInteractiveDeferredLoadingClientConfig(rules, &overrides, os.Stdin)
}

func getClientset() (*kubernetes.Clientset, error) {
	loader := getClientConfigLoader()

	config, err := loader.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (server *Server) getDeBotOpsSecret() (*corev1.Secret, error) {
	return server.kubernetesClientset.
		CoreV1().
		Secrets("default").
		Get(context.TODO(), core.DeBotOpsSecretName, metav1.GetOptions{})
}

type DeBotOpsClientset struct {
	Listener dynamic.NamespaceableResourceInterface
}

func NewDeBotOpsClientset() (*DeBotOpsClientset, error) {
	loader := getClientConfigLoader()

	config, err := loader.ClientConfig()
	if err != nil {
		return nil, err
	}

	dynamicInterface, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	listenerInterface := v1alpha1.GroupVersion.WithResource("listeners")

	clientset := &DeBotOpsClientset{
		Listener: dynamicInterface.Resource(listenerInterface),
	}

	return clientset, nil
}
