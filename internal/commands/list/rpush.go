package list

import (
	"fmt"

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
	if len(args) != 2 {
		return WrongNumberOfArgumentsError("rpush")
	}

	key := args[0].String()
	items := []any{}

	for _, arg := range args[1:] {
		items = append(items, arg.String())
	}

	entry, exists := c.storage.Get(key)
	if !exists {
		list := store.NewList()
		list.Append(items)
		c.storage.Set(key, list)
		return resp.NewInteger(fmt.Sprintf("%d", len(args)-1))
	}

	list, isList := entry.(*store.List)
	if !isList {
		return WrongTypeOperationError()
	}

	list.Append(items)
	c.storage.Set(key, list)

	return resp.NewInteger(fmt.Sprintf("%d", list.Size()))
}

func WrongNumberOfArgumentsError(command string) resp.Value {
	return resp.NewSimpleError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", command))
}

func WrongTypeOperationError() resp.Value {
	return resp.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value")
}
