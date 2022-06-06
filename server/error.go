package server

import (
	"github.com/gofiber/fiber/v2"
)

func fiberErrorHandler(ctx *fiber.Ctx, fiberError *fiber.Error) error {
	return ctx.
		Status(fiberError.Code).
		JSON(&Response{
			Status: StatusError,
			Data:   fiberError.Message,
		})
}

func defaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if fiberError, ok := err.(*fiber.Error); ok {
		return fiberErrorHandler(ctx, fiberError)
	}
	return ctx.
		Status(fiber.StatusInternalServerError).
		JSON(&Response{
			Status: StatusError,
		})
}

func jwtErrorHandler(ctx *fiber.Ctx, _ error) error {
	return ctx.
		Status(fiber.StatusUnauthorized).
		JSON(&Response{
			Status: StatusError,
			Data:   "Invalid or expired JWT",
		})
}
