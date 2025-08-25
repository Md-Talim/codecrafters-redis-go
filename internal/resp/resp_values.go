package resp

import "fmt"

type Value interface {
	String() string
	Serialize() []byte
}

type Array struct {
	items []Value
}

func NewArray(items []Value) *Array {
	return &Array{items}
}

func (a *Array) Items() []Value {
	return a.items
}

func (a *Array) String() string {
	result := ""
	for i, item := range a.items {
		if i > 0 {
			result += " "
		}
		result += item.String()
	}
	return result
}

func (a *Array) Serialize() []byte {
	result := fmt.Appendf(nil, "*%d%s", len(a.items), CRLF)
	for _, item := range a.items {
		result = append(result, item.Serialize()...)
	}
	return result
}

type SimpleString struct {
	value string
}

func NewSimpleString(value string) *SimpleString {
	return &SimpleString{value}
}

func (s *SimpleString) String() string {
	return s.value
}

func (s *SimpleString) Serialize() []byte {
	return fmt.Appendf(nil, "+%s%s", s.value, CRLF)
}

type SimpleError struct {
	value string
}

func NewSimpleError(value string) *SimpleError {
	return &SimpleError{value}
}

func (s *SimpleError) String() string {
	return s.value
}

func (s *SimpleError) Serialize() []byte {
	return fmt.Appendf(nil, "-%s%s", s.value, CRLF)
}

type BulkString struct {
	value string
}

func NewBulkString(value string) *BulkString {
	return &BulkString{value}
}

func NewNullBulkString() *BulkString {
	return &BulkString{"null"}
}

func (s *BulkString) Type() string { return "BulkString " }

func (s *BulkString) String() string { return s.value }

func (s *BulkString) Serialize() []byte {
	if s.value == "null" {
		return []byte("$-1\r\n")
	}
	return fmt.Appendf(nil, "$%d%s%s%s", len(s.value), CRLF, s.value, CRLF)
}
