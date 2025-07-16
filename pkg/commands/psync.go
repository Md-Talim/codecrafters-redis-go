package commands

import (
	"fmt"

	"github.com/md-talim/codecrafters-redis-go/internal/rdb"
	"github.com/md-talim/codecrafters-redis-go/pkg/replication"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type PsyncCommand struct {
	replication *replication.Info
}

func NewPsyncCommand(replication *replication.Info) *PsyncCommand {
	return &PsyncCommand{replication}
}

func (c *PsyncCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 2 {
		return WrongNumberOfArgumentsError("psync")
	}

	if args[0].Type != resp.BulkString || args[1].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	replID := args[0].Bulk
	offset := args[1].Bulk

	if replID == "?" && offset == "-1" {
		response := fmt.Sprintf("FULLRESYNC %s 0", c.replication.MasterReplID())
		return &resp.Value{
			Type:    resp.RDBFile,
			String:  response,
			RDBData: rdb.EmptyRDBFile(),
		}
	}

	response := fmt.Sprintf("FULLRESYNC %s 0", c.replication.MasterReplID())
	return resp.NewSimpleString(response)
}

func (c *PsyncCommand) Name() string {
	return "PSYNC"
}
