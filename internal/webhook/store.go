package webhook

import (
	"fmt"
	"sync"
	"time"
)

// Simple in-memory store (we'll add Redis later)
type Store struct {
	endpoints map[string]*Endpoint
	webhooks  map[string][]*Webhook
	mu        sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		endpoints: make(map[string]*Endpoint),
		webhooks:  make(map[string][]*Webhook),
	}
}

func (s *Store) CreateEndpoint(id string) (*Endpoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	endpoint := &Endpoint{
		ID:        id,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24h expiry
	}
	s.endpoints[id] = endpoint
	s.webhooks[id] = []*Webhook{}

	return endpoint, nil
}

func (s *Store) GetEndpoint(id string) (*Endpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	endpoint, exists := s.endpoints[id]
	if !exists {
		return nil, fmt.Errorf("endpoint not found")
	}
	return endpoint, nil
}

func (s *Store) StoreWebhook(endpointID string, webhook *Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.endpoints[endpointID]; !exists {
		return fmt.Errorf("endpoint not found")
	}

	s.webhooks[endpointID] = append(s.webhooks[endpointID], webhook)

	// Keep only last 100 webhooks per endpoint
	if len(s.webhooks[endpointID]) > 100 {
		s.webhooks[endpointID] = s.webhooks[endpointID][1:]
	}

	return nil
}

func (s *Store) GetWebhooks(endpointID string) ([]*Webhook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	webhooks, exists := s.webhooks[endpointID]
	if !exists {
		return nil, fmt.Errorf("endpoint not found")
	}
	return webhooks, nil
}
