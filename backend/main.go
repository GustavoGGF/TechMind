package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Static("/assets", "./build/browser/assets")

    app.Static("/", "./build/browser")


    local := app.Group("/")      // /api

    local.Get("/api", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "name": "Grame",
            "age": 20,
          })
    })      
  
    log.Fatal(app.Listen(":3000"))
}