// Package auth emite y valida tokens JWT (HS256).
package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// LocalsToken es la clave donde el middleware guarda el token validado,
// para que los handlers puedan reenviarlo a la API de Node.
const LocalsToken = "token"

// GenerateToken crea un JWT firmado con HS256 para el usuario indicado.
func GenerateToken(secret []byte, subject string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// ValidateToken verifica firma, algoritmo y expiración.
func ValidateToken(secret []byte, tokenString string) error {
	_, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (any, error) { return secret, nil },
		// Solo aceptamos HS256: evita ataques de confusión de algoritmo ("alg": "none").
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	return err
}

// Middleware exige un header "Authorization: Bearer <token>" válido.
func Middleware(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := bearerToken(c.Get(fiber.HeaderAuthorization))
		if err == nil {
			err = ValidateToken(secret, token)
		}
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token ausente, inválido o expirado",
			})
		}
		c.Locals(LocalsToken, token)
		return c.Next()
	}
}

func bearerToken(header string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) || len(header) == len(prefix) {
		return "", errors.New("header Authorization ausente o mal formado")
	}
	return strings.TrimPrefix(header, prefix), nil
}
