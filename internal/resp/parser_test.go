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

	if _, ok := value.(*SimpleString); !ok {
		t.Errorf("Expected SimpleString")
	}

	if value.String() != "PONG" {
		t.Errorf("Expected 'PONG', got %q", value.String())
	}
}

func TestParseBulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	parser := NewParser(strings.NewReader(input))

	value, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if _, ok := value.(*BulkString); !ok {
		t.Error("Expected BulkString")
	}

	if value.String() != "hello" {
		t.Errorf("Expected 'hello'")
	}
}

func TestParseArray(t *testing.T) {
	input := "*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"
	parser := NewParser(strings.NewReader(input))

	value, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	array, ok := value.(*Array)
	if !ok {
		t.Errorf("Expected Array")
	}

	items := array.Items()
	if len(items) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(items))
	}

	if items[0].String() != "ECHO" {
		t.Errorf("Expected 'ECHO', got %q", items[0].String())
	}

	if items[1].String() != "hello" {
		t.Errorf("Expected 'hello', got %q", items[1].String())
	}
}
