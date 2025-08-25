package resp

import "testing"

func TestSerializeSimpleString(t *testing.T) {
	value := NewSimpleString("PONG")
	expected := "+PONG\r\n"
	actual := value.Serialize()

	if string(actual) != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSerializeBulkString(t *testing.T) {
	value := NewBulkString("hello")
	expected := "$5\r\nhello\r\n"
	actual := value.Serialize()

	if string(actual) != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSerializeSimpleError(t *testing.T) {
	value := NewSimpleError("ERR unknown command")
	expected := "-ERR unknown command\r\n"
	actual := value.Serialize()

	if string(actual) != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}
