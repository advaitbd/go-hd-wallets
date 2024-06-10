package main

import (
	"go-keystone/mod/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
}

func main() {
	app := fiber.New()
	setupMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/fingerprint", handlers.DecodeURHandler)
	app.Post("/transaction", handlers.GenerateSignRequestHandler)
	app.Post("/sign", handlers.SignTransactionHandler)

	log.Fatal(app.Listen(":8000"))
}
