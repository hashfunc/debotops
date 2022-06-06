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

	token, err := jwt.ParseWithClaims(refreshTokenCookie, &auth.Claims{}, server.keyFunc)
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

	return ctx.JSON(
		&Response{
			Status: StatusOK,
			Data: fiber.Map{
				"token": access.String(),
			},
		},
	)
}
