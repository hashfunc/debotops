package server

import (
	"golang.org/x/crypto/bcrypt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/hashfunc/debotops/pkg/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (server *Server) login(ctx *fiber.Ctx) error {
	payload := &LoginRequest{}

	if err := ctx.BodyParser(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if server.root == nil {
		log.Println("`debotops-secret` is invalid")
		return fiber.ErrInternalServerError
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(server.root.PasswordHash),
		[]byte(payload.Password+server.root.SecretKey),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	refresh, err := auth.NewAuth(server.root.SecretKey, auth.KindRefresh)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Value:    refresh.String(),
		Expires:  refresh.Expires,
		HTTPOnly: true,
	})

	access, err := auth.NewAuth(server.root.SecretKey, auth.KindAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	data := fiber.Map{
		"token": access.String(),
	}

	return ctx.JSON(data)
}
