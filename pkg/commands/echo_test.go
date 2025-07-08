package commands

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func TestEchoCommand(t *testing.T) {
	cmd := &EchoCommand{}
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "grape"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.BulkString {
		t.Errorf("Expected BulkString, got %v", result.Type)
	}

	if result.Bulk != "grape" {
		t.Errorf("Expected 'grape', got %q", result.Bulk)
	}
}

func TestEchoCommandNoArgs(t *testing.T) {
	cmd := &EchoCommand{}
	args := []resp.Value{}

	result := cmd.Execute(args)

	if result.Type != resp.SimpleError {
		t.Errorf("Expected Error, got %v", result.Type)
	}
}
