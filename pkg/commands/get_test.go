package commands

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func TestGetCommandExists(t *testing.T) {
	memoryStorage := storage.NewInMemory()
	memoryStorage.Set("foo", "bar")

	getCmd := NewGetCommand(memoryStorage)
	getCmdArgs := []resp.Value{
		{Type: resp.BulkString, Bulk: "foo"},
	}
	result := getCmd.Execute(getCmdArgs)

	if result.Type != resp.BulkString {
		t.Errorf("Expected BulkString, got %q", result.Type)
	}

	if result.Bulk != "bar" {
		t.Errorf("Expected 'bar', got %q", result.String)
	}
}

func TestGetCommandNotExists(t *testing.T) {
	memoryStorage := storage.NewInMemory()
	cmd := NewGetCommand(memoryStorage)
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "nonexistent"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.BulkString {
		t.Errorf("Expected BulkString (null), got %v", result.Type)
	}

	if result.Serialize() != "$-1\r\n" {
		t.Errorf("Expected null bulk string, got %q", result.Serialize())
	}
}
