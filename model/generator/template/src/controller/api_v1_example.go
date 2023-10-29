package controller

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/gofiber/fiber/v2/middleware/session"
)

func newApiV1Controller(r fiber.Router, store *Store) {
	controller := ApiV1Controller{NewController(r, store)}

	r.Post("activate/admin", controller.activateAdmin)
	r.Post("login", controller.login)
}
