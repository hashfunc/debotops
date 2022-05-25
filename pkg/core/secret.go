package core

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"math/rand"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Root struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	SecretKey    string `json:"secret_key"`
}

const DeBotOpsSecretName = "debotops-secret"

func NewDeBotOpsSecret(namespace string) ([]*corev1.Secret, error) {
	var secrets []*corev1.Secret

	username := generateRandomString(16)
	password := generateRandomString(32)
	secretKey := generateRandomString(64)

	passwordHash, err := encryptPassword(password, secretKey)
	if err != nil {
		return nil, err
	}

	root := Root{
		Username:     username,
		PasswordHash: passwordHash,
		SecretKey:    secretKey,
	}

	rootJson, err := json.Marshal(&root)
	if err != nil {
		return nil, err
	}

	initialSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeBotOpsSecretName + "-initial",
			Namespace: namespace,
		},
		StringData: map[string]string{
			"password": password,
		},
	}
	secrets = append(secrets, initialSecret)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeBotOpsSecretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"root": rootJson,
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

	return string(hashBytes), nil
}
