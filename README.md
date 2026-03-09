# mailgate-go

Go client for the [mailgate](https://github.com/techpartners-asia/mailgate) HTTP mail API.

## Install

```bash
go get github.com/techpartners-asia/mailgate-go/client
```

## Usage

```go
import "github.com/techpartners-asia/mailgate-go/client"

c := client.New("http://localhost:8025", "your-api-key")

// Send email
err := c.Send(ctx, client.SendRequest{
    To:       []string{"user@example.com"},
    Subject:  "Welcome",
    BodyText: "Plain text body",
    BodyHTML: "<p>HTML body</p>",
})

// With attachments
err = c.Send(ctx, client.SendRequest{
    To:       []string{"user@example.com"},
    Subject:  "Invoice",
    BodyText: "See attachment.",
    Attachments: []client.Attachment{
        {Filename: "invoice.pdf", ContentType: "application/pdf", Data: pdfBytes},
    },
})

// List send logs
logs, count, err := c.Logs(ctx, client.LogsFilter{Status: "sent", Limit: 20, Offset: 0})

// Health check (no API key)
health, err := c.Health(ctx)
```

API errors (4xx/5xx) are returned as `*client.APIError` with `StatusCode` and `Message`. Use `client.NewWithClient(baseURL, apiKey, httpClient)` for custom timeouts or TLS.

## Requirements

- Go 1.21+
- A running [mailgate](https://github.com/techpartners-asia/mailgate) server
