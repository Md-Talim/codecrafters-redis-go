package list

import (
	"fmt"

	"github.com/md-talim/codecrafters-redis-go/internal/commands/core"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type RPushCommand struct {
	storage store.Storage
}

func NewRPushCommand(storage store.Storage) *RPushCommand {
	return &RPushCommand{storage}
}

func (c *RPushCommand) Execute(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return core.WrongNumberOfArgumentsError("rpush")
	}

	key := args[0].String()
	items := args[1:]

	entry, exists := c.storage.Get(key)
	if !exists {
		list := store.NewList()
		list.Append(items)
		c.storage.Set(key, list)
		return resp.NewInteger(fmt.Sprintf("%d", len(args)-1))
	}

	list, isList := entry.(*store.List)
	if !isList {
		return core.WrongTypeOperationError()
	}

	list.Append(items)
	c.storage.Set(key, list)

	return resp.NewInteger(fmt.Sprintf("%d", list.Size()))
}
