package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"reto-interseguro/api-go/internal/auth"
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

// mustEnv exige que la variable exista: falla al arrancar en lugar de
// funcionar con un secreto inseguro por defecto.
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("falta la variable de entorno %s", key)
	}
	return v
}

func main() {
	statsURL := getEnv("STATS_API_URL", "http://localhost:3000")
	port := getEnv("PORT", "8080")
	jwtSecret := []byte(mustEnv("JWT_SECRET"))
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET debe tener al menos 32 caracteres")
	}

	authHandler := handlers.NewAuthHandler(mustEnv("AUTH_USERNAME"), mustEnv("AUTH_PASSWORD"), jwtSecret, time.Hour)
	qrHandler := handlers.NewQRHandler(client.NewStatsClient(statsURL))

	app := fiber.New()

	// Rutas públicas
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "api-go"})
	})
	app.Post("/api/auth/login", authHandler.Login)

	// Rutas protegidas con JWT
	app.Post("/api/qr", auth.Middleware(jwtSecret), qrHandler.Handle)

	log.Fatal(app.Listen(":" + port))
}
