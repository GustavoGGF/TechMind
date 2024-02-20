package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Post("/credentials", func(c *fiber.Ctx) error {
        var data map[string]string

        if err := c.BodyParser(&data); err != nil {
            return err
        }

        return c.JSON(([]string{}))
    })

    app.Listen(":3000")
}