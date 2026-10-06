// Package handlers contiene los endpoints HTTP de la API.
package handlers

import (
	"github.com/gofiber/fiber/v2"

	"reto-interseguro/api-go/internal/matrix"
)

// QRRequest es el cuerpo esperado: {"matrix": [[1,2],[3,4]]}
type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// QRResponse devuelve las dos matrices de la factorización.
type QRResponse struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// ErrorResponse es el formato común de error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// QR maneja POST /api/qr.
func QR(c *fiber.Ctx) error {
	var req QRRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "JSON inválido: se espera {\"matrix\": [[números]]}",
		})
	}
	if err := matrix.Validate(req.Matrix); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	q, r := matrix.QR(req.Matrix)
	return c.JSON(QRResponse{Q: q, R: r})
}
