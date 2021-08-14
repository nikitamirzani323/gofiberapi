package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/gofiberapi/model"
)

type ClientToken struct {
	Token string `json:"token"`
}
type ClientInit struct {
	Client_Company string `json:"client_company"`
}

func Fetch_token(c *fiber.Ctx) error {
	client := new(ClientToken)

	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":          fiber.StatusOK,
		"token":           client.Token,
		"member_username": "developer",
		"member_company":  "MMD",
		"member_credit":   5000000,
	})
}
func FetchAll_pasaran(c *fiber.Ctx) error {
	client := new(ClientInit)

	if err := c.BodyParser(client); err != nil {
		return err
	}

	result, err := model.FetchAll_MclientPasaran(client.Client_Company)

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
