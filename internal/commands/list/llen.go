package list

import (
	"fmt"

	"github.com/md-talim/codecrafters-redis-go/internal/commands/core"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type LLenCommand struct {
	storage store.Storage
}

func NewLLenCommand(storage store.Storage) *LLenCommand {
	return &LLenCommand{storage}
}

func (c *LLenCommand) Execute(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return core.WrongNumberOfArgumentsError("llen")
	}

	key := args[0].String()
	entry, exists := c.storage.Get(key)
	if !exists {
		return resp.NewInteger("0")
	}

	list, isList := entry.(*store.List)
	if !isList {
		return core.WrongTypeOperationError()
	}

	return resp.NewInteger(fmt.Sprintf("%d", list.Size()))
}
