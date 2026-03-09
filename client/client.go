// Package client provides a Go SDK for the mailgate HTTP API.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client talks to a mailgate server.
type Client struct {
	BaseURL    string       // e.g. "http://localhost:8025"
	APIKey     string       // X-API-Key header for /send and /logs
	HTTPClient *http.Client // optional; default http.DefaultClient
}

// New returns a Client. baseURL should not have a trailing slash.
func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
}

// NewWithClient returns a Client using the given http.Client.
func NewWithClient(baseURL, apiKey string, hc *http.Client) *Client {
	c := New(baseURL, apiKey)
	c.HTTPClient = hc
	return c
}

func (c *Client) url(path string) string {
	return c.BaseURL + path
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, auth bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.url(path), body)
	if err != nil {
		return nil, err
	}
	if auth {
		req.Header.Set("X-API-Key", c.APIKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func decodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var errBody struct {
			OK    bool   `json:"ok"`
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return &APIError{StatusCode: resp.StatusCode, Message: errBody.Error}
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// APIError is returned when the server responds with 4xx or 5xx.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("mailgate: %d %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("mailgate: HTTP %d", e.StatusCode)
}
