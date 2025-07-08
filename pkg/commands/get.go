package commands

import (
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type GetCommand struct {
	storage storage.Storage
}

func NewGetCommand(storage storage.Storage) *GetCommand {
	return &GetCommand{storage}
}

func (g *GetCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("get")
	}

	if args[0].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	key := args[0].Bulk
	value, exists := g.storage.Get(key)
	if !exists {
		return resp.NewNullBulkString()
	}

	return resp.NewBulkString(value)
}

func (g *GetCommand) Name() string {
	return "GET"
}
