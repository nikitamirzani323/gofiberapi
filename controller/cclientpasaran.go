package controller

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/gofiberapi/config"
	"github.com/nikitamirzani323/gofiberapi/model"
)

type ClientToken struct {
	Token string `json:"token"`
}
type ClientInit struct {
	Client_Company string `json:"client_company"`
}

var ctx = context.Background()

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
	conf := config.GetConfigRedis()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.DB_HOST,
		Password: conf.DB_PASSWORD,
		DB:       conf.DB_NAME,
	})

	resultredis, err := rdb.Get(ctx, "listpasaran_"+client.Client_Company).Result()

	if err == redis.Nil {
		result, err := model.FetchAll_MclientPasaran(client.Client_Company)

		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		json, _ := json.Marshal(result)
		log.Println("mysql")
		err = rdb.Set(ctx, "listpasaran_"+client.Client_Company, json, 0).Err()
		if err != nil {
			panic(err)
		}
		return c.JSON(result)
	} else {
		log.Println("cache")
		return c.SendString(resultredis)
	}
}
