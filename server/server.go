package server

import (
	"net/http"

	"k8s.io/client-go/kubernetes"
)

type Server struct {
	mux       *http.ServeMux
	clientset *kubernetes.Clientset
}

func NewServer() (*Server, error) {
	clientset, err := getClientset()
	if err != nil {
		return nil, err
	}

	server := &Server{
		mux:       http.NewServeMux(),
		clientset: clientset,
	}

	server.mux.HandleFunc("/login", server.login)
	server.mux.HandleFunc("/refresh", server.refresh)

	return server, nil
}

func (server *Server) Start() error {
	return http.ListenAndServe(":8386", server.mux)
}
