package Response

import "github.com/gofiber/fiber/v2"

func RespErr(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"asdf": "asdf",
	})
}

func RespOK(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"asdf": "asdf",
	})
}

func RespCreated(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"asdf": "asdf",
	})
}
