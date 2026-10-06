// Package client contiene el cliente HTTP hacia la API de estadísticas (Node).
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Stats refleja la respuesta de POST /api/stats de la API de Node.
type Stats struct {
	Max              float64 `json:"max"`
	Min              float64 `json:"min"`
	Average          float64 `json:"average"`
	Sum              float64 `json:"sum"`
	Count            int     `json:"count"`
	AnyDiagonal      bool    `json:"anyDiagonal"`
	DiagonalByMatrix []bool  `json:"diagonalByMatrix"`
}

// StatsClient llama a la API de Node.
type StatsClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewStatsClient crea el cliente con un timeout para no quedar colgados
// si la API de Node no responde.
func NewStatsClient(baseURL string) *StatsClient {
	return &StatsClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Compute envía las matrices a Node y devuelve las estadísticas.
func (c *StatsClient) Compute(ctx context.Context, token string, matrices ...[][]float64) (*Stats, error) {
	body, err := json.Marshal(map[string]any{"matrices": matrices})
	if err != nil {
		return nil, fmt.Errorf("serializando matrices: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/stats", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("llamando a la API de estadísticas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("la API de estadísticas respondió %d", resp.StatusCode)
	}

	var stats Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("leyendo respuesta de estadísticas: %w", err)
	}
	return &stats, nil
}
