package commands

import (
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type SetCommand struct {
	storage storage.Storage
}

func NewSetCommand(storage storage.Storage) *SetCommand {
	return &SetCommand{storage}
}

func (s *SetCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 2 {
		return resp.NewSimpleError("ERR wrong number of arguments for 'set' command")
	}

	if args[0].Type != resp.BulkString || args[1].Type != resp.BulkString {
		return resp.NewSimpleError("Err invalid argument type")
	}

	key := args[0].Bulk
	value := args[1].Bulk

	err := s.storage.Set(key, value)
	if err != nil {
		return resp.NewSimpleError("ERR failed to set key")
	}

	return resp.NewSimpleString("OK")
}

func (s *SetCommand) Name() string {
	return "SET"
}
