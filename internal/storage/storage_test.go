package storage

import (
	"testing"
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
)

func TestNewInMemory(t *testing.T) {
	storage := NewInMemory()
	defer storage.Close()

	// Test basic operations
	err := storage.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, exists := storage.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Get failed: expected value1, got %s", value)
	}
}

func TestNewWithConfig(t *testing.T) {
	// Test with empty config (no RDB)
	cfg := &config.Config{
		Dir:        "",
		DBFilename: "",
		Port:       "6379",
	}

	storage := New(cfg)

	// Should work as normal memory storage
	storage.Set("test", "value")
	value, exists := storage.Get("test")

	if !exists || value != "value" {
		t.Errorf("Expected test=value, got %s", value)
	}
}

func TestRDBLoading(t *testing.T) {
	// Test with RDB config but non-existent file
	cfg := &config.Config{
		Dir:        "/tmp/nonexistent",
		DBFilename: "nonexistent.rdb",
		Port:       "6379",
	}

	storage := New(cfg)

	// Should still work (empty storage)
	keys := storage.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected empty storage, got %d keys", len(keys))
	}
}

func TestExpiry(t *testing.T) {
	storage := NewInMemory()
	defer storage.Close()

	// Set key with very short expiry
	storage.SetWithExpiry("temp", "value", 10*time.Millisecond)

	// Should exist immediately
	value, exists := storage.Get("temp")
	if !exists || value != "value" {
		t.Error("Key should exist immediately after setting")
	}

	// Wait for expiry
	time.Sleep(20 * time.Millisecond)

	// Should be expired
	_, exists = storage.Get("temp")
	if exists {
		t.Error("Key should be expired")
	}
}
