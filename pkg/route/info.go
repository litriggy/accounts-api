package route

import (
	"accounts/api/app/controller"

	"github.com/gofiber/fiber/v2"
)

func InfoRouter(r fiber.Router) {
	r.Get("/services", controller.GetTotalServices)
}
