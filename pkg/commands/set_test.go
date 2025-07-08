package commands

import (
	"testing"
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func TestSetCommand(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	cmd := NewSetCommand(memoryStorage)
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "foo"},
		{Type: resp.BulkString, Bulk: "bar"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.SimpleString {
		t.Errorf("Expected SimpleString, got %q", result.Type)
	}
	if result.String != "OK" {
		t.Errorf("Expected 'PONG', got %q", result.String)
	}

	value, exists := memoryStorage.Get("foo")
	if !exists {
		t.Errorf("Expected key 'foo' to exist")
	}
	if value != "bar" {
		t.Errorf("Expected 'bar', got %q", value)
	}
}

func TestSetCommandInvalidArgs(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	cmd := NewSetCommand(memoryStorage)
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "foo"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.SimpleError {
		t.Errorf("Expected SimpleError, got %v", result.Type)
	}
}

func TestSetCommandWithPX(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	cmd := NewSetCommand(memoryStorage)

	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "foo"},
		{Type: resp.BulkString, Bulk: "bar"},
		{Type: resp.BulkString, Bulk: "PX"},
		{Type: resp.BulkString, Bulk: "100"}, // 100ms
	}

	result := cmd.Execute(args)

	if result.Type != resp.SimpleString {
		t.Errorf("Expected SimpleString, got %v", result.Type)
	}

	if result.String != "OK" {
		t.Errorf("Expected 'OK', got %q", result.String)
	}

	// Verify key exists immediately
	value, exists := memoryStorage.Get("foo")
	if !exists {
		t.Error("Expected key 'foo' to exist")
	}
	if value != "bar" {
		t.Errorf("Expected 'bar', got %q", value)
	}

	// Wait for expiry
	time.Sleep(150 * time.Millisecond)

	// Verify key expired
	_, exists = memoryStorage.Get("foo")
	if exists {
		t.Error("Expected key 'foo' to be expired")
	}
}
