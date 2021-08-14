package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/gofiberapi/model"
)

func FetchAll_pasaran(c *fiber.Ctx) error {
	result, err := model.FetchAll_MclientPasaran()

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	return c.JSON(result)
}
