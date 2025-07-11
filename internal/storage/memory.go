package storage

import (
	"sync"
	"time"
)

type StorageItem struct {
	Value     string
	ExpriesAt *time.Time
}

func newStorageItem(value string, expiresAt *time.Time) *StorageItem {
	return &StorageItem{value, expiresAt}
}

type MemoryStorage struct {
	data map[string]*StorageItem
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		data: make(map[string]*StorageItem),
	}

	go storage.cleanupExpiredKeys()

	return storage
}

func (m *MemoryStorage) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = newStorageItem(value, nil)
	return nil
}

func (m *MemoryStorage) SetWithExpiry(key, value string, expiry time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiryTime := time.Now().Add(expiry)
	m.data[key] = newStorageItem(value, &expiryTime)
	return nil
}

func (m *MemoryStorage) Get(key string) (string, bool) {
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

func (m *MemoryStorage) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryStorage) Keys() []string {
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

func (m *MemoryStorage) cleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()

		for key, item := range m.data {
			if item.ExpriesAt != nil && time.Now().After(*item.ExpriesAt) {
				delete(m.data, key)
			}
		}

		m.mu.Unlock()
	}
}
