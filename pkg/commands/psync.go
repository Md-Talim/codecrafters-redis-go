package commands

import (
	"fmt"

	"github.com/md-talim/codecrafters-redis-go/internal/replica"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type PsyncCommand struct {
	replication *replica.Info
}

func NewPsyncCommand(replication *replica.Info) *PsyncCommand {
	return &PsyncCommand{replication}
}

func (c *PsyncCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 2 {
		return WrongNumberOfArgumentsError("psync")
	}

	if args[0].Type != resp.BulkString || args[1].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	response := fmt.Sprintf("FULLRESYNC %s 0", c.replication.MasterReplID())
	return resp.NewSimpleString(response)
}

func (c *PsyncCommand) Name() string {
	return "PSYNC"
}
