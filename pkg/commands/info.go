package commands

import (
	"strings"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type InfoCommand struct {
	config *config.Config
}

func NewInfoCommand(config *config.Config) *InfoCommand {
	return &InfoCommand{config}
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
	response := "role:master"
	return resp.NewBulkString(response)
}

func (i *InfoCommand) Name() string {
	return "INFO"
}
