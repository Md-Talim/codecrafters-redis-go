package commands

import "github.com/md-talim/codecrafters-redis-go/pkg/resp"

type ReplConfCommand struct{}

func NewReplConfCommand() *ReplConfCommand {
	return &ReplConfCommand{}
}

func (r *ReplConfCommand) Execute(args []resp.Value) *resp.Value {
	// Ignore the args and send OK response
	return resp.NewSimpleString("OK")
}

func (r *ReplConfCommand) Name() string {
	return "REPLCONF"
}
