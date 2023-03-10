package route

import (
	"accounts/api/app/controller"
	"accounts/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r fiber.Router) {
	r.Get("/wallet", middleware.CheckSession)
	r.Post("/wallet/:type", middleware.CheckSession, controller.AddWallet)
	r.Delete("/wallet", middleware.CheckSession)
	r.Get("/balance", middleware.CheckSession, controller.GetBalance)
	r.Get("/services", middleware.CheckSession, controller.GetUserServices)
	r.Post("/service/:serviceId", middleware.CheckSession, controller.AddService)
	r.Post("/secondpass", middleware.CheckSession, controller.CreateSecondPassword)
}
