package middleware

import (
	"accounts/api/pkg/config"
	"accounts/api/platform/memcached"

	"github.com/gofiber/fiber/v2"
)

func CheckSession(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	sessionKey := headers["Authorization"]
	if config.AppCfg().Stage == "test" {
		c.Locals("userID", sessionKey)
		c.Set("Authorization", "new")
		return c.Next()
	}
	if sessionKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "No accessKey provided",
		})
	}
	userID, newSessionKey, err := memcached.GetSession(sessionKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"msg":    "Error on finding Session Key",
			"data":   err.Error(),
		})
	}
	c.Set("Authorization", newSessionKey)
	c.Locals("userID", userID)
	return c.Next()
}
