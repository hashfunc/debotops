package core

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"math/rand"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	RootName   string
	RootSecret string
}

const DeBotOpsSecretName = "debotops-secret"

func NewDeBotOpsSecret(namespace string) ([]*corev1.Secret, error) {
	var secrets []*corev1.Secret

	rootName := generateRandomString(16)
	rootSecret := generateRandomString(64)
	rootPassword := generateRandomString(32)

	rootPasswordHash, err := encryptPassword(rootPassword, rootSecret)
	if err != nil {
		return nil, err
	}

	initialSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeBotOpsSecretName + "-initial",
			Namespace: namespace,
		},
		StringData: map[string]string{
			"RootPassword": rootPassword,
		},
	}
	secrets = append(secrets, initialSecret)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeBotOpsSecretName,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"RootName":         rootName,
			"RootSecret":       rootSecret,
			"RootPasswordHash": rootPasswordHash,
		},
	}
	secrets = append(secrets, secret)

	return secrets, nil
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	buffer := make([]byte, length)
	for index := range buffer {
		buffer[index] = characters[rand.Intn(len(characters))]
	}

	return string(buffer)
}

func encryptPassword(password, secret string) (string, error) {
	plain := []byte(password + secret)

	hashBytes, err := bcrypt.GenerateFromPassword(plain, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	base64HashString := base64.StdEncoding.EncodeToString(hashBytes)
	return base64HashString, nil
}
