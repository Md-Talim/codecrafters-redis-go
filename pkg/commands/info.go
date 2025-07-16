package commands

import (
	"strings"

	"github.com/md-talim/codecrafters-redis-go/pkg/replication"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type InfoCommand struct {
	replication *replication.Info
}

func NewInfoCommand(replication *replication.Info) *InfoCommand {
	return &InfoCommand{replication}
}

func (i *InfoCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) > 1 {
		return WrongNumberOfArgumentsError("info")
	}

	var section string
	if len(args) == 1 {
		if args[0].Type != resp.BulkString {
			return InvalidArgumentTypeError()
		}
		section = strings.ToLower(args[0].Bulk)
	} else {
		section = "all"
	}

	switch section {
	case "replication":
		return i.getReplicationInfo()
	case "all":
		return i.getReplicationInfo()
	default:
		return resp.NewBulkString("")
	}
}

func (i *InfoCommand) getReplicationInfo() *resp.Value {
	response := i.replication.InfoString()
	return resp.NewBulkString(response)
}

func (i *InfoCommand) Name() string {
	return "INFO"
}
