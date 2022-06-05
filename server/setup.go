package server

import jwtmiddleware "github.com/gofiber/jwt/v3"

func (server *Server) setup() {
	server.setRoute()
}

func (server *Server) setRoute() {
	server.fiber.Post("/login", server.login)
	server.fiber.Post("/refresh", server.refresh)

	server.fiber.Use(jwtmiddleware.New(
		jwtmiddleware.Config{KeyFunc: server.keyFunc},
	))
}
