package resp

type Type int

const (
	SimpleString Type = iota
	SimpleError
	BulkString
	Array
)

type Value struct {
	Type   Type
	String string
	Bulk   string
	Array  []Value
}
