package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"reto-interseguro/api-go/internal/handlers"
)

func main() {
	app := fiber.New()

	// Endpoint de salud, útil luego para Docker y la nube
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "api-go"})
	})

	app.Post("/api/qr", handlers.QR)

	log.Fatal(app.Listen(":8080"))
}
