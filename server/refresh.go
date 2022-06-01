package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/hashfunc/debotops/pkg/auth"
)

func (server *Server) refresh(ctx *fiber.Ctx) error {
	refreshTokenCookie := ctx.Cookies("refresh-token")
	if refreshTokenCookie == "" {
		return fiber.ErrUnauthorized
	}

	token, err := jwt.ParseWithClaims(refreshTokenCookie, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.root.SecretKey), nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := token.Claims.(*auth.Claims)

	if err := claims.Valid(); err != nil {
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
