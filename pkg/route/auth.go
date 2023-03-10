package route

import (
	"accounts/api/app/controller"
	"accounts/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r fiber.Router) {
	r.Post("/signin/:type/:version", controller.SignIn)
	r.Get("/check", middleware.CheckSession, controller.CheckSession)
}
