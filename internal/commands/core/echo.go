package core

import "github.com/md-talim/codecrafters-redis-go/internal/resp"

type EchoCommand struct{}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

func (e *EchoCommand) Execute(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("echo")
	}

	return resp.NewBulkString(args[0].String())
}

func (p *EchoCommand) Name() string {
	return "ECHO"
}
