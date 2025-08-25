package commands

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

func TestGetCommandExists(t *testing.T) {
	memoryStorage := store.NewInMemory()
	memoryStorage.Set("foo", "bar")

	getCmd := NewGetCommand(memoryStorage)
	getCmdArgs := []resp.Value{
		resp.NewBulkString("foo"),
	}

	result := getCmd.Execute(getCmdArgs)

	if _, isBulkString := result.(*resp.BulkString); !isBulkString {
		t.Errorf("Expected BulkString")
	}

	if result.String() != "bar" {
		t.Errorf("Expected 'bar', got %q", result.String())
	}
}

func TestGetCommandNotExists(t *testing.T) {
	memoryStorage := store.NewInMemory()
	cmd := NewGetCommand(memoryStorage)
	args := []resp.Value{
		resp.NewBulkString("nonexistent"),
	}

	result := cmd.Execute(args)

	if _, isBulkString := result.(*resp.BulkString); !isBulkString {
		t.Errorf("Expected BulkString (null)")
	}

	if result.String() != "null" {
		t.Errorf("Expected null bulk string, got %q", result.Serialize())
	}
}
