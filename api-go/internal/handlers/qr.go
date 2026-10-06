// Package handlers contiene los endpoints HTTP de la API.
package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"reto-interseguro/api-go/internal/auth"
	"reto-interseguro/api-go/internal/client"
	"reto-interseguro/api-go/internal/matrix"
)

// StatsProvider abstrae la API de estadísticas. Usar una interfaz permite
// reemplazarla por un "fake" en los tests (como una interfaz inyectada en C#).
type StatsProvider interface {
	Compute(ctx context.Context, token string, matrices ...[][]float64) (*client.Stats, error)
}

// QRRequest es el cuerpo esperado: {"matrix": [[1,2],[3,4]]}
type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// QRResponse devuelve la factorización y las estadísticas calculadas por Node.
type QRResponse struct {
	Q     [][]float64   `json:"q"`
	R     [][]float64   `json:"r"`
	Stats *client.Stats `json:"stats"`
}

// ErrorResponse es el formato común de error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// QRHandler agrupa las dependencias del endpoint.
type QRHandler struct {
	stats StatsProvider
}

// NewQRHandler recibe sus dependencias por parámetro (inyección de dependencias).
func NewQRHandler(stats StatsProvider) *QRHandler {
	return &QRHandler{stats: stats}
}

// Handle maneja POST /api/qr: valida, factoriza, pide estadísticas a Node
// y devuelve todo junto.
func (h *QRHandler) Handle(c *fiber.Ctx) error {
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

	// Reenviamos el mismo token del usuario: Node también lo valida.
	token, _ := c.Locals(auth.LocalsToken).(string)
	stats, err := h.stats.Compute(c.UserContext(), token, q, r)
	if err != nil {
		log.Printf("error en API de estadísticas: %v", err)
		return c.Status(fiber.StatusBadGateway).JSON(ErrorResponse{
			Error: "no se pudieron calcular las estadísticas",
		})
	}

	return c.JSON(QRResponse{Q: q, R: r, Stats: stats})
}
