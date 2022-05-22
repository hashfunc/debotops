package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/hashfunc/debotops/pkg/core"
)

func GetSecret(client client.Client, ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	secret := new(corev1.Secret)
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}

	err := client.Get(ctx, key, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func GetDeBotOpsSecret(client client.Client, ctx context.Context, namespace string) (*corev1.Secret, error) {
	return GetSecret(client, ctx, namespace, core.DeBotOpsSecretName)
}
