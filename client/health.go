package client

import (
	"context"
	"net/http"
)

// HealthResponse is the response from GET /health. The endpoint reports
// liveness only; the aggregate-stats body was removed server-side (per-merchant
// stats now live behind the admin API).
type HealthResponse struct {
	OK   bool   `json:"ok"`
	Code string `json:"code,omitempty"`
}

// Health calls GET /health (no API key required). Returns server liveness.
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
