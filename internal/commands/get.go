package commands

import (
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type GetCommand struct {
	storage store.Storage
}

func NewGetCommand(storage store.Storage) *GetCommand {
	return &GetCommand{storage}
}

func (g *GetCommand) Execute(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("get")
	}

	key := args[0].String()
	value, exists := g.storage.Get(key)
	if !exists {
		return resp.NewNullBulkString()
	}

	return resp.NewBulkString(value)
}

func (g *GetCommand) Name() string {
	return "GET"
}
