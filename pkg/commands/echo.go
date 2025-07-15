package commands

import "github.com/md-talim/codecrafters-redis-go/pkg/resp"

type EchoCommand struct{}

func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

func (e *EchoCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("echo")
	}

	if args[0].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	return resp.NewBulkString(args[0].Bulk)
}

func (p *EchoCommand) Name() string {
	return "ECHO"
}
