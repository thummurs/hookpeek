package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/thummurs/hookpeek/internal/webhook"
)

type API struct {
	store *webhook.Store
}

func NewAPI(store *webhook.Store) *API {
	return &API{store: store}
}

// Generate random ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Create new webhook endpoint
func (a *API) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	id := generateID()
	endpoint, err := a.store.CreateEndpoint(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"endpoint": endpoint,
		"url":      "http://" + r.Host + "/w/" + id,
	})
}

// Get endpoint details
func (a *API) GetEndpoint(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/endpoints/"):]

	endpoint, err := a.store.GetEndpoint(id)
	if err != nil {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	webhooks, _ := a.store.GetWebhooks(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"endpoint": endpoint,
		"webhooks": webhooks,
		"count":    len(webhooks),
	})
}

// Capture incoming webhook
func (a *API) CaptureWebhook(w http.ResponseWriter, r *http.Request) {
	endpointID := r.URL.Path[len("/w/"):]

	// Verify endpoint exists
	_, err := a.store.GetEndpoint(endpointID)
	if err != nil {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	// Read body
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Extract headers
	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}

	// Extract query params
	query := make(map[string]string)
	for key, values := range r.URL.Query() {
		query[key] = values[0]
	}

	// Create webhook record
	webhook := &webhook.Webhook{
		ID:        generateID(),
		Method:    r.Method,
		Path:      r.URL.Path,
		Headers:   headers,
		Body:      string(body),
		Query:     query,
		IP:        r.RemoteAddr,
		Timestamp: time.Now(),
	}

	// Store webhook
	if err := a.store.StoreWebhook(endpointID, webhook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook captured"))
}

// Health check
func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
