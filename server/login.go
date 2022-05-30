package server

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"

	"github.com/hashfunc/debotops/pkg/auth"
	"github.com/hashfunc/debotops/pkg/core"
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

	secret, err := server.getDeBotOpsSecret()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	secretData := secret.Data

	var root core.Root
	if err := json.Unmarshal(secretData["root"], &root); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(root.PasswordHash),
		[]byte(payload.Password+root.SecretKey),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	refresh, err := auth.NewAuth(root.SecretKey, auth.KindRefresh)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh-token",
		Value:    refresh.String(),
		Expires:  refresh.Expires,
		HTTPOnly: true,
	})

	access, err := auth.NewAuth(root.SecretKey, auth.KindAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	data := fiber.Map{
		"token": access.String(),
	}

	return ctx.JSON(data)
}
