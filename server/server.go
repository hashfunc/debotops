package server

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"k8s.io/client-go/kubernetes"

	"github.com/hashfunc/debotops/pkg/core"
	"github.com/hashfunc/debotops/server/pipe"
)

type Server struct {
	root                *core.Root
	fiber               *fiber.App
	kubernetesClientset *kubernetes.Clientset
	debotopsClientset   *DeBotOpsClientset
	waitGroup           *sync.WaitGroup
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

	fiberConfig := fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          defaultErrorHandler,
	}

	server := &Server{
		fiber:               fiber.New(fiberConfig),
		kubernetesClientset: kubernetesClientset,
		debotopsClientset:   debotopsClientset,
		waitGroup:           &sync.WaitGroup{},
	}

	server.setup()

	return server, nil
}

func (server *Server) Start() {
	server.runPipe()
	server.runServer()
	server.waitGroup.Wait()
}

func (server *Server) runServer() {
	server.waitGroup.Add(1)

	go func() {
		err := server.fiber.Listen(":8386")
		log.Fatal(err)
	}()
}

func (server *Server) runPipe() {
	server.waitGroup.Add(2)

	namespaceForSecret := core.GetNamespaceForSecret()

	secretPipe := pipe.NewSecretPipe(namespaceForSecret, server.kubernetesClientset)
	go secretPipe.Run()

	go func() {
		for {
			select {

			case secret := <-secretPipe.Channel:
				root := new(core.Root)
				if err := json.Unmarshal(secret.Data["root"], root); err != nil {
					server.root = nil
				}
				server.root = root
			}
		}
	}()
}

func (server *Server) keyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(server.root.SecretKey), nil
}
