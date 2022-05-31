package server

import (
	"github.com/gofiber/fiber/v2"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	fiber               *fiber.App
	kubernetesClientset *kubernetes.Clientset
	debotopsClientset   *DeBotOpsClientset
}

func NewServer() (*Server, error) {
	kubernetesClientset, err := getClientset()
	if err != nil {
		return nil, err
	}

	debotopsClientset, err := NewDeBotOpsClientset()
	if err != nil {
		return nil, err
	}

	server := &Server{
		fiber:               fiber.New(),
		kubernetesClientset: kubernetesClientset,
		debotopsClientset:   debotopsClientset,
	}

	server.fiber.Post("/login", server.login)
	server.fiber.Post("/refresh", server.refresh)

	return server, nil
}

func (server *Server) Start() error {
	return server.fiber.Listen(":8386")
}
