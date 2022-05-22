package core

import (
	"math/rand"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	RootName   string
	RootSecret string
}

const DeBotOpsSecretName = "debotops-secret"

func NewDeBotOpsSecret(namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeBotOpsSecretName,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"RootName":   generateRandomString(16),
			"RootSecret": generateRandomString(64),
		},
	}
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	buffer := make([]byte, length)
	for index := range buffer {
		buffer[index] = characters[rand.Intn(len(characters))]
	}

	return string(buffer)
}
