package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (server *Server) listListeners(ctx *fiber.Ctx) error {
	listeners, err := server.debotopsClientset.
		Listener.List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(listeners)
}

func (server *Server) getListener(ctx *fiber.Ctx) error {
	namespace := ctx.Params("namespace")
	name := ctx.Params("name")

	listener, err := server.debotopsClientset.
		Listener.Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(listener)
}
