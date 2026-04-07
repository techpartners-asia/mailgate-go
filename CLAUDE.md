# Mailgate Go SDK

Go client library for the Mailgate HTTP API.

## Purpose

Provides a typed Go client for sending email, checking health, and querying logs via the Mailgate gateway. Used by backend services to integrate with Mailgate without raw HTTP calls.

## Structure

```
client/
  client.go       # Client struct, HTTP helpers, APIError type
  send.go         # Send() — POST /send
  health.go       # Health() — GET /health
  logs.go         # Logs() — GET /logs
  client_test.go  # Tests for all methods
```

## Usage

```go
import "github.com/techpartners-asia/mailgate-go/client"

c := client.New("http://localhost:8025", "merchant-api-key")

// Send email
err := c.Send(ctx, client.SendRequest{
    To:       []string{"user@example.com"},
    Subject:  "Hello",
    BodyText: "Plain text",
})

// Health (no auth)
health, err := c.Health(ctx)

// Logs (scoped to merchant)
logs, count, err := c.Logs(ctx, client.LogsFilter{Limit: 20})
```

## Key Types

- `Client` — holds base URL, API key, HTTP client
- `SendRequest` — to, cc, bcc, subject, body_text, body_html, attachments
- `Attachment` — filename, content_type, data (raw bytes, auto base64-encoded)
- `HealthResponse` — ok, stats
- `LogEntry` — id, sent_at, merchant_name, to, subject, status, error, duration_ms
- `APIError` — status code + message for 4xx/5xx responses

## Notes

- This is a **separate Go module** (`github.com/techpartners-asia/mailgate-go`), included as a git submodule in the main mailgate repo
- Only covers merchant-facing endpoints (send, health, logs) — no admin API
- API key is the merchant's key, passed via `X-API-Key` header
- Errors from the server are returned as `*client.APIError`

## Commands

```bash
go test ./...    # Run tests
```
