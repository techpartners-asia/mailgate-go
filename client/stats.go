package client

import (
	"context"
	"net/http"
	"time"
)

// CategoryStats is one category's counts in GET /stats.
type CategoryStats struct {
	Category string `json:"category"` // "" for mail sent without a category
	Sent     int    `json:"sent"`
	Failed   int    `json:"failed"`
	Cap      int    `json:"cap"` // daily cap for this category; 0 = none
}

// Stats is the response from GET /stats: the authenticated merchant's counts
// since the start of the UTC+8 day.
type Stats struct {
	Since      time.Time       `json:"since"`
	DailyCap   int             `json:"daily_cap"` // 0 = unlimited
	Sent       int             `json:"sent"`
	Failed     int             `json:"failed"`
	Categories []CategoryStats `json:"categories"`
}

// Stats returns today's sent/failed counts per category (UTC+8 day).
func (c *Client) Stats(ctx context.Context) (*Stats, error) {
	resp, err := c.do(ctx, http.MethodGet, "/stats", nil, true)
	if err != nil {
		return nil, err
	}
	var env struct {
		Data Stats `json:"data"`
	}
	if err := decodeJSON(resp, &env); err != nil {
		return nil, err
	}
	return &env.Data, nil
}
