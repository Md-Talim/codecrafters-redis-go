package core

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
)

func TestPingCommandNoArgs(t *testing.T) {
	cmd := &PingCommand{}
	args := []resp.Value{}

	result := cmd.Execute(args)

	if _, isString := result.(*resp.SimpleString); !isString {
		t.Errorf("Expected SimpleString")
	}

	if result.String() != "PONG" {
		t.Errorf("Expected 'PONG', got %q", result.String())
	}
}

func TestPingCommandWithArg(t *testing.T) {
	cmd := &PingCommand{}
	args := []resp.Value{
		resp.NewBulkString("hello"),
	}

	result := cmd.Execute(args)

	if _, isBulk := result.(*resp.BulkString); !isBulk {
		t.Errorf("Expected SimpleString")
	}

	if result.String() != "hello" {
		t.Errorf("Expected 'hello', got %q", result.String())
	}
}
