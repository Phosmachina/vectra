package controller

import (
	. "Vectra/src/view/go"
	"github.com/gofiber/fiber/v2"
	. "github.com/gofiber/fiber/v2/middleware/session"
	"io"
)

func newViewController(r fiber.Router, store *Store) {
	controller := ViewController{NewController(r, store)}

	r.Get("/", controller.root)
	r.Get("/init", controller.init)
	r.Get("/login", controller.login)
	r.Get("/sign", controller.sign)
}

func (c ViewController) root(ctx *fiber.Ctx) error {
	return HandleView(ctx, c, func(buf io.Writer, userId string) error {
		Jade_index(NewGlobalCtx("Index", userId), buf)
		return nil
	})
}

func (c ViewController) init(ctx *fiber.Ctx) error {
	return HandleView(ctx, c, func(buf io.Writer, userId string) error {
		Jade_init(NewGlobalCtx("Initialization", userId), buf)
		return nil
	})
}

func (c ViewController) login(ctx *fiber.Ctx) error {
	return HandleView(ctx, c, func(buf io.Writer, userId string) error {
		if userId != "" {
			return ctx.Redirect("/", fiber.StatusPreconditionRequired)
		}
		Jade_login(NewGlobalCtx("Login", userId), buf)
		return nil
	})
}

func (c ViewController) sign(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
