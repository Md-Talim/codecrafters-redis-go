package resp

import "fmt"

func (v *Value) Serialize() string {
	switch v.Type {
	case SimpleString:
		return fmt.Sprintf("+%s%s", v.String, CRLF)
	case BulkString:
		if v.String == "null" {
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d%s%s%s", len(v.Bulk), CRLF, v.Bulk, CRLF)
	case SimpleError:
		return fmt.Sprintf("-%s%s", v.String, CRLF)
	case Array:
		result := fmt.Sprintf("*%d%s", len(v.Array), CRLF)
		for _, item := range v.Array {
			result += item.Serialize()
		}
		return result
	case RDBFile:
		fullresync := NewSimpleString(v.String).Serialize()
		rdbFile := fmt.Sprintf("$%d%s%s", len(v.RDBData), CRLF, string(v.RDBData))
		return fullresync + rdbFile
	default:
		return ""
	}
}

func NewSimpleString(value string) *Value {
	return &Value{Type: SimpleString, String: value}
}

func NewBulkString(data string) *Value {
	return &Value{Type: BulkString, Bulk: data}
}

func NewSimpleError(message string) *Value {
	return &Value{Type: SimpleError, String: message}
}

func NewNullBulkString() *Value {
	return &Value{Type: BulkString, Bulk: "", String: "null"}
}

func NewArray(values []Value) *Value {
	return &Value{Type: Array, Array: values}
}
