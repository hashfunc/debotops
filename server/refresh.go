package server

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/hashfunc/debotops/pkg/auth"
	"github.com/hashfunc/debotops/pkg/core"
)

func (server *Server) refresh(ctx *fiber.Ctx) error {
	refreshTokenCookie := ctx.Cookies("refresh-token")
	if refreshTokenCookie == "" {
		return fiber.ErrUnauthorized
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

	token, err := jwt.ParseWithClaims(refreshTokenCookie, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(root.SecretKey), nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := token.Claims.(*auth.Claims)

	if err := claims.Valid(); err != nil {
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
