package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"reto-interseguro/api-go/internal/client"
	"reto-interseguro/api-go/internal/handlers"
)

// getEnv devuelve la variable de entorno o un valor por defecto.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	statsURL := getEnv("STATS_API_URL", "http://localhost:3000")
	port := getEnv("PORT", "8080")

	qrHandler := handlers.NewQRHandler(client.NewStatsClient(statsURL))

	app := fiber.New()

	// Endpoint de salud, útil luego para Docker y la nube
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "api-go"})
	})

	app.Post("/api/qr", qrHandler.Handle)

	log.Fatal(app.Listen(":" + port))
}
