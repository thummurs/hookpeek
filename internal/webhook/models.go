package webhook

import (
	"time"
)

type Webhook struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Query     map[string]string `json:"query"`
	IP        string            `json:"ip"`
	Timestamp time.Time         `json:"timestamp"`
}

type Endpoint struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}