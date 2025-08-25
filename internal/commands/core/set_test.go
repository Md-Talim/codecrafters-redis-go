package core

import (
	"testing"
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

func TestSetCommand(t *testing.T) {
	memoryStorage := store.NewInMemory()
	cmd := NewSetCommand(memoryStorage)
	args := []resp.Value{
		resp.NewBulkString("foo"),
		resp.NewBulkString("bar"),
	}

	result := cmd.Execute(args)

	if _, isString := result.(*resp.SimpleString); !isString {
		t.Errorf("Expected SimpleString")
	}
	if result.String() != "OK" {
		t.Errorf("Expected 'PONG', got %q", result.String())
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
	memoryStorage := store.NewInMemory()
	cmd := NewSetCommand(memoryStorage)
	args := []resp.Value{
		resp.NewBulkString("foo"),
	}

	result := cmd.Execute(args)

	if _, isError := result.(*resp.SimpleError); !isError {
		t.Errorf("Expected SimpleError")
	}
}

func TestSetCommandWithPX(t *testing.T) {
	memoryStorage := store.NewInMemory()
	cmd := NewSetCommand(memoryStorage)

	args := []resp.Value{
		resp.NewBulkString("foo"),
		resp.NewBulkString("bar"),
		resp.NewBulkString("PX"),
		resp.NewBulkString("100"),
	}

	result := cmd.Execute(args)

	if _, isString := result.(*resp.SimpleString); !isString {
		t.Errorf("Expected SimpleString")
	}

	if result.String() != "OK" {
		t.Errorf("Expected 'OK', got %q", result.String())
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
