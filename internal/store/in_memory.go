package store

import (
	"sync"
	"time"
)

type Item struct {
	Value     string
	ExpriesAt *time.Time
}

type InMemory struct {
	data   map[string]*Item
	mu     sync.RWMutex
	closer chan struct{}
}

func newItem(value string, expiresAt *time.Time) *Item {
	return &Item{value, expiresAt}
}

func NewInMemory() *InMemory {
	storage := &InMemory{
		data:   make(map[string]*Item),
		closer: make(chan struct{}),
	}

	go storage.cleanupExpiredKeys()

	return storage
}

func (m *InMemory) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = newItem(value, nil)
	return nil
}

func (m *InMemory) SetWithExpiry(key, value string, expiry time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiryTime := time.Now().Add(expiry)
	m.data[key] = newItem(value, &expiryTime)
	return nil
}

func (m *InMemory) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return "", false
	}

	// Check if expired and delete
	if item.ExpriesAt != nil && time.Now().After(*item.ExpriesAt) {
		m.mu.RUnlock()
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()
		m.mu.RLock()
		return "", false
	}

	return item.Value, true
}

func (m *InMemory) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *InMemory) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.data))
	now := time.Now()

	for key, item := range m.data {
		if item.ExpriesAt != nil && now.After(*item.ExpriesAt) {
			continue
		}
		keys = append(keys, key)
	}

	return keys
}

func (m *InMemory) Close() {
	close(m.closer)
}

func (m *InMemory) cleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			for key, item := range m.data {
				if item.ExpriesAt != nil && now.After(*item.ExpriesAt) {
					delete(m.data, key)
				}
			}
			m.mu.Unlock()

		case <-m.closer:
			return
		}
	}
}
