package resp

import (
	"strings"
	"testing"
)

func TestParseSimpleString(t *testing.T) {
	input := "+PONG\r\n"
	parser := NewParser(strings.NewReader(input))

	value, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if value.Type != SimpleString {
		t.Errorf("Expected SimpleString, got %v", value.Type)
	}

	if value.String != "PONG" {
		t.Errorf("Expected 'PONG', got %q", value.String)
	}
}

func TestParseBulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	parser := NewParser(strings.NewReader(input))

	value, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if value.Type != BulkString {
		t.Errorf("Expected BulkString, got %v", value.Type)
	}

	if value.Bulk != "hello" {
		t.Errorf("Expected 'hello', got %q", value.Bulk)
	}
}

func TestParseArray(t *testing.T) {
	input := "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"
	parser := NewParser(strings.NewReader(input))

	value, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if value.Type != Array {
		t.Errorf("Expected Array, got %v", value.Type)
	}

	if len(value.Array) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(value.Array))
	}

	if value.Array[0].Bulk != "ECHO" {
		t.Errorf("Expected 'ECHO', got %q", value.Array[0].Bulk)
	}

	if value.Array[1].Bulk != "hello" {
		t.Errorf("Expected 'hello', got %q", value.Array[1].Bulk)
	}
}
