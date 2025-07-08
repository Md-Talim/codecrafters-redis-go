package commands

import "github.com/md-talim/codecrafters-redis-go/pkg/resp"

type PingCommand struct{}

func (p *PingCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) == 0 {
		return resp.NewSimpleString("PONG")
	}

	if len(args) == 1 && args[0].Type == resp.BulkString {
		return resp.NewBulkString(args[0].Bulk)
	}

	return resp.NewSimpleError("ERR wrong number of arguments for 'ping' command")
}

func (p *PingCommand) Name() string {
	return "PING"
}
