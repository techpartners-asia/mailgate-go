package client

import (
	"context"
	"net/http"
)

// HealthResponse is the response from GET /health.
type HealthResponse struct {
	OK       bool   `json:"ok"`
	Provider string `json:"provider"`
	Stats    *Stats `json:"stats,omitempty"`
}

// Stats holds aggregate send stats (optional in health response).
type Stats struct {
	TotalSent     int64 `json:"total_sent"`
	TotalFailed   int64 `json:"total_failed"`
	Last24hSent   int64 `json:"last_24h_sent"`
	Last24hFailed int64 `json:"last_24h_failed"`
}

// Health calls GET /health (no API key required). Returns server liveness and optional stats.
func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	resp, err := c.do(ctx, http.MethodGet, "/health", nil, false)
	if err != nil {
		return HealthResponse{}, err
	}
	var out HealthResponse
	if err := decodeJSON(resp, &out); err != nil {
		return HealthResponse{}, err
	}
	return out, nil
}
