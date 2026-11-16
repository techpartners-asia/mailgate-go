package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stats" || r.Method != http.MethodGet || r.Header.Get("X-API-Key") != "key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"ok":true,"code":"OK","data":{"since":"2026-11-16T00:00:00+08:00","daily_cap":0,"sent":469,"failed":3,
"categories":[{"category":"ebarimt","sent":412,"failed":3,"cap":0},{"category":"auth","sent":57,"failed":0,"cap":500}]}}`))
	}))
	defer server.Close()

	got, err := New(server.URL, "key").Stats(context.Background())
	if err != nil {
		t.Fatalf("Stats: %v", err)
	}
	if got.Sent != 469 || got.Failed != 3 || got.DailyCap != 0 || len(got.Categories) != 2 {
		t.Fatalf("Stats = %+v", got)
	}
	if got.Since.Format("2006-01-02T15:04:05Z07:00") != "2026-11-16T00:00:00+08:00" {
		t.Errorf("Since = %v", got.Since)
	}
	if c := got.Categories[1]; c.Category != "auth" || c.Sent != 57 || c.Cap != 500 {
		t.Errorf("Categories[1] = %+v", c)
	}
}

func TestStats_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "code": "UNAUTHORIZED", "error": "unauthorized"})
	}))
	defer server.Close()

	got, err := New(server.URL, "bad").Stats(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || got != nil {
		t.Fatalf("Stats: got %v, err %v, want *APIError", got, err)
	}
	if apiErr.StatusCode != 401 || apiErr.Code != "UNAUTHORIZED" {
		t.Errorf("apiErr = %+v", apiErr)
	}
}
