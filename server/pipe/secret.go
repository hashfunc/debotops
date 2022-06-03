package pipe

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type SecretPipe struct {
	namespace string
	clientset *kubernetes.Clientset
	Channel   chan *corev1.Secret
}

func NewSecretPipe(namespace string, clientset *kubernetes.Clientset) *SecretPipe {
	return &SecretPipe{
		namespace: namespace,
		clientset: clientset,
		Channel:   make(chan *corev1.Secret),
	}
}

func (pipe *SecretPipe) Run() {
	options := v1.ListOptions{
		FieldSelector: "metadata.name=debotops-secret",
	}

	watcher, err := pipe.clientset.
		CoreV1().
		Secrets(pipe.namespace).
		Watch(context.TODO(), options)

	if err != nil {
		log.Fatal(err)
	}

	for event := range watcher.ResultChan() {
		switch event.Type {

		case watch.Added, watch.Modified:
			secret, ok := event.Object.(*corev1.Secret)
			if !ok {
				continue
			}

			pipe.Channel <- secret

		case watch.Deleted:
			pipe.Channel <- nil
		}
	}
}
