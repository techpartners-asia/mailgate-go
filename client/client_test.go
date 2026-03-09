package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	c := New("http://localhost:8025", "secret")
	if c.BaseURL != "http://localhost:8025" || c.APIKey != "secret" {
		t.Errorf("New: BaseURL=%q APIKey=%q", c.BaseURL, c.APIKey)
	}
	if c.HTTPClient != http.DefaultClient {
		t.Error("New: expected DefaultClient")
	}
}

func TestAPIError_Error(t *testing.T) {
	e := &APIError{StatusCode: 401, Message: "unauthorized"}
	if e.Error() != "mailgate: 401 unauthorized" {
		t.Errorf("Error() = %q", e.Error())
	}
	e2 := &APIError{StatusCode: 500, Message: ""}
	if e2.Error() != "mailgate: HTTP 500" {
		t.Errorf("Error() = %q", e2.Error())
	}
}

func TestSend_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/send" || r.Method != http.MethodPost {
			t.Errorf("unexpected path/method: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "key" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "unauthorized"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}))
	defer server.Close()

	c := New(server.URL, "key")
	err := c.Send(context.Background(), SendRequest{
		To:       []string{"a@test.com"},
		Subject:  "Hi",
		BodyText: "Hello",
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
}

func TestSend_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "invalid email"})
	}))
	defer server.Close()

	c := New(server.URL, "key")
	err := c.Send(context.Background(), SendRequest{
		To:       []string{"a@test.com"},
		Subject:  "Hi",
		BodyText: "Hello",
	})
	if err == nil {
		t.Fatal("Send: expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Send: expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 || apiErr.Message != "invalid email" {
		t.Errorf("APIError: %d %q", apiErr.StatusCode, apiErr.Message)
	}
}

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("path = %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":       true,
			"provider": "smtp",
			"stats": map[string]int64{
				"total_sent": 10, "total_failed": 1,
				"last_24h_sent": 2, "last_24h_failed": 0,
			},
		})
	}))
	defer server.Close()

	c := New(server.URL, "key")
	resp, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if !resp.OK || resp.Provider != "smtp" {
		t.Errorf("Health: ok=%v provider=%q", resp.OK, resp.Provider)
	}
	if resp.Stats == nil || resp.Stats.TotalSent != 10 {
		t.Errorf("Health: stats = %+v", resp.Stats)
	}
}

func TestLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/logs" || r.Header.Get("X-API-Key") != "key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logs": []map[string]interface{}{
				{
					"id": "1", "sent_at": "2025-03-09T12:00:00Z", "provider": "smtp",
					"to": []string{"a@test.com"}, "subject": "S", "status": "sent",
					"duration_ms": int64(100),
				},
			},
			"count": 1,
		})
	}))
	defer server.Close()

	c := New(server.URL, "key")
	logs, count, err := c.Logs(context.Background(), LogsFilter{Limit: 10})
	if err != nil {
		t.Fatalf("Logs: %v", err)
	}
	if count != 1 || len(logs) != 1 {
		t.Errorf("Logs: count=%d len(logs)=%d", count, len(logs))
	}
	if logs[0].ID != "1" || logs[0].Subject != "S" || logs[0].Status != "sent" {
		t.Errorf("Logs[0]: %+v", logs[0])
	}
}

func TestSend_ValidationError(t *testing.T) {
	c := New("http://localhost:8025", "key")
	err := c.Send(context.Background(), SendRequest{
		To:      []string{"a@test.com"},
		Subject: "Hi",
		// no body_text or body_html
	})
	if err == nil {
		t.Fatal("Send: expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		if err.Error() != "mailgate: 0 at least one of body_text or body_html required" {
			t.Errorf("error = %v", err)
		}
		return
	}
	_ = apiErr
}
