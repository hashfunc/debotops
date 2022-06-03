package server

import (
	"os"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hashfunc/debotops/api/v1alpha1"
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
