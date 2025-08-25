package commands

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
)

func TestEchoCommand(t *testing.T) {
	cmd := &EchoCommand{}
	args := []resp.Value{
		resp.NewBulkString("grape"),
	}

	result := cmd.Execute(args)

	value, isBulkString := result.(*resp.BulkString)
	if !isBulkString {
		t.Errorf("Expected BulkString")
	}

	if value.String() != "grape" {
		t.Errorf("Expected 'grape', got %q", value.String())
	}
}

func TestEchoCommandNoArgs(t *testing.T) {
	cmd := &EchoCommand{}
	args := []resp.Value{}

	result := cmd.Execute(args)

	if _, isError := result.(*resp.SimpleError); !isError {
		t.Errorf("Expected BulkString")
	}
}
