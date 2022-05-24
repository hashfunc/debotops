package server

import (
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	server := &Server{
		mux: http.NewServeMux(),
	}

	server.mux.HandleFunc("/", server.greeting)

	return server
}

func (server *Server) Start() error {
	return http.ListenAndServe(":8386", server.mux)
}

func (server *Server) greeting(writer http.ResponseWriter, _ *http.Request) {
	message := "Greeting from API server"
	if _, err := writer.Write([]byte(message)); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
