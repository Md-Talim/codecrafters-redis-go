package core

import (
	"fmt"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
)

func InvalidArgumentTypeError() resp.Value {
	return resp.NewSimpleError("ERR invalid argument type")
}

func WrongNumberOfArgumentsError(command string) resp.Value {
	return resp.NewSimpleError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", command))
}

func SyntaxError() resp.Value {
	return resp.NewSimpleError("ERR syntax error")
}

func UnknownCommandError(command string) resp.Value {
	return resp.NewSimpleError(fmt.Sprintf("ERR unknown command '%s'", command))
}

func InvalidExpireTimeError(command string) resp.Value {
	return resp.NewSimpleError(fmt.Sprintf("ERR invalid expire time in '%s' command", command))
}

func InvalidCommandFormatError() resp.Value {
	return resp.NewSimpleError("ERR invalid command format")
}

func ValueNotIntegerError() resp.Value {
	return resp.NewSimpleError("ERR value is not an integer or out of range")
}

func WrongTypeOperationError() resp.Value {
	return resp.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value")
}
