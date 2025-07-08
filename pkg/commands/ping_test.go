package commands

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func TestPingCommandNoArgs(t *testing.T) {
	cmd := &PingCommand{}
	args := []resp.Value{}

	result := cmd.Execute(args)

	if result.Type != resp.SimpleString {
		t.Errorf("Expected SimpleString, got %q", result.Type)
	}

	if result.String != "PONG" {
		t.Errorf("Expected 'PONG', got %q", result.String)
	}
}

func TestPingCommandWithArg(t *testing.T) {
	cmd := &PingCommand{}
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "hello"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.BulkString {
		t.Errorf("Expected SimpleString, got %q", result.Type)
	}

	if result.Bulk != "hello" {
		t.Errorf("Expected 'hello', got %q", result.Bulk)
	}
}
