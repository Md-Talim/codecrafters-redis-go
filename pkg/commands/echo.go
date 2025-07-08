package commands

import "github.com/md-talim/codecrafters-redis-go/pkg/resp"

type EchoCommand struct{}

func (e *EchoCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 1 {
		return resp.NewSimpleError("ERR wrong number of arguments for 'echo' command")
	}

	if args[0].Type != resp.BulkString {
		return resp.NewSimpleError("ERR invalid argument type")
	}

	return resp.NewBulkString(args[0].Bulk)
}

func (p *EchoCommand) Name() string {
	return "ECHO"
}
