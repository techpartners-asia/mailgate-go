package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// SendRequest is the payload for POST /send.
type SendRequest struct {
	To          []string     `json:"to"`
	CC          []string     `json:"cc,omitempty"`
	BCC         []string     `json:"bcc,omitempty"`
	Subject     string       `json:"subject"`
	BodyText    string       `json:"body_text,omitempty"`
	BodyHTML    string       `json:"body_html,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment is one attachment in a send request (data is base64-encoded in JSON).
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"-"` // sent as base64 in JSON
}

type sendAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
}

type sendRequestJSON struct {
	To          []string         `json:"to"`
	CC          []string         `json:"cc,omitempty"`
	BCC         []string         `json:"bcc,omitempty"`
	Subject     string           `json:"subject"`
	BodyText    string           `json:"body_text,omitempty"`
	BodyHTML    string           `json:"body_html,omitempty"`
	Attachments []sendAttachment `json:"attachments,omitempty"`
}

// Send sends an email via the mailgate API. At least one of BodyText or BodyHTML must be set.
func (c *Client) Send(ctx context.Context, req SendRequest) error {
	if req.BodyText == "" && req.BodyHTML == "" {
		return &APIError{StatusCode: 0, Message: "at least one of body_text or body_html required"}
	}
	payload := sendRequestJSON{
		To:       req.To,
		CC:       req.CC,
		BCC:      req.BCC,
		Subject:  req.Subject,
		BodyText: req.BodyText,
		BodyHTML: req.BodyHTML,
	}
	for _, a := range req.Attachments {
		payload.Attachments = append(payload.Attachments, sendAttachment{
			Filename:    a.Filename,
			ContentType: a.ContentType,
			Data:        base64.StdEncoding.EncodeToString(a.Data),
		})
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := c.do(ctx, http.MethodPost, "/send", bytes.NewReader(body), true)
	if err != nil {
		return err
	}
	var out struct {
		OK bool `json:"ok"`
	}
	return decodeJSON(resp, &out)
}
