package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// LogsFilter is the query for GET /logs.
type LogsFilter struct {
	Status string // "sent" or "failed"; empty = all
	Limit  int    // max results (server caps at 500); 0 = server default
	Offset int
}

// LogEntry is one entry from GET /logs.
type LogEntry struct {
	ID           string    `json:"id"`
	SentAt       time.Time `json:"sent_at"`
	MerchantName string    `json:"merchant_name"`
	To           []string  `json:"to"`
	Subject      string    `json:"subject"`
	Status       string    `json:"status"`
	Error        string    `json:"error,omitempty"`
	DurationMS   int64     `json:"duration_ms"`
}

// LogsResponse is the response from GET /logs.
type LogsResponse struct {
	Logs  []LogEntry `json:"logs"`
	Count int       `json:"count"`
}

// Logs returns send log entries. count is the number of entries in this page.
func (c *Client) Logs(ctx context.Context, filter LogsFilter) (logs []LogEntry, count int, err error) {
	q := url.Values{}
	if filter.Status != "" {
		q.Set("status", filter.Status)
	}
	if filter.Limit > 0 {
		q.Set("limit", strconv.Itoa(filter.Limit))
	}
	if filter.Offset > 0 {
		q.Set("offset", strconv.Itoa(filter.Offset))
	}
	path := "/logs"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	resp, err := c.do(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, 0, err
	}
	var out LogsResponse
	if err := decodeJSON(resp, &out); err != nil {
		return nil, 0, err
	}
	return out.Logs, out.Count, nil
}
