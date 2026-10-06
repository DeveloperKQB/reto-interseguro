package handlers

import (
	"crypto/subtle"
	"time"

	"github.com/gofiber/fiber/v2"

	"reto-interseguro/api-go/internal/auth"
)

// LoginRequest: {"username": "...", "password": "..."}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse devuelve el token y su duración en segundos.
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}

// AuthHandler valida credenciales de demostración (configuradas por
// variables de entorno) y emite un JWT.
type AuthHandler struct {
	username, password string
	secret             []byte
	ttl                time.Duration
}

// NewAuthHandler recibe sus dependencias por parámetro.
func NewAuthHandler(username, password string, secret []byte, ttl time.Duration) *AuthHandler {
	return &AuthHandler{username: username, password: password, secret: secret, ttl: ttl}
}

// Login maneja POST /api/auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "JSON inválido"})
	}

	// Comparación en tiempo constante: evita ataques de temporización.
	userOK := subtle.ConstantTimeCompare([]byte(req.Username), []byte(h.username)) == 1
	passOK := subtle.ConstantTimeCompare([]byte(req.Password), []byte(h.password)) == 1
	if !userOK || !passOK {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Error: "credenciales inválidas"})
	}

	token, err := auth.GenerateToken(h.secret, req.Username, h.ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "no se pudo generar el token"})
	}
	return c.JSON(LoginResponse{Token: token, ExpiresIn: int(h.ttl.Seconds())})
}
