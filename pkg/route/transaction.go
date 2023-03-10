package route

import (
	"accounts/api/app/controller"
	"accounts/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(r fiber.Router) {
	r.Post("/transfer", middleware.CheckSession, controller.TransferBalance)
	r.Get("/history", middleware.CheckSession, controller.TransferHistory)
}
