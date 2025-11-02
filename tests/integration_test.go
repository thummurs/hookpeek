package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thummurs/hookpeek/internal/api"
	"github.com/thummurs/hookpeek/internal/webhook"
)

func TestCreateEndpoint(t *testing.T) {
	store := webhook.NewStore()
	apiHandler := api.NewAPI(store)

	req := httptest.NewRequest("POST", "/api/endpoints", nil)
	w := httptest.NewRecorder()

	apiHandler.CreateEndpoint(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if result["url"] == nil {
		t.Error("Expected URL in response")
	}
}

func TestCaptureWebhook(t *testing.T) {
	store := webhook.NewStore()
	apiHandler := api.NewAPI(store)

	// Create endpoint first
	endpoint, _ := store.CreateEndpoint("test123")

	// Send webhook
	body := bytes.NewBufferString(`{"test":"data"}`)
	req := httptest.NewRequest("POST", "/w/"+endpoint.ID, body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiHandler.CaptureWebhook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Verify webhook was stored
	webhooks, _ := store.GetWebhooks(endpoint.ID)
	if len(webhooks) != 1 {
		t.Errorf("Expected 1 webhook, got %d", len(webhooks))
	}
}
