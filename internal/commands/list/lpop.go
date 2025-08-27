package list

import (
	"strconv"

	"github.com/md-talim/codecrafters-redis-go/internal/commands/core"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type LPopCommand struct {
	storage store.Storage
}

func NewLPopCommand(storage store.Storage) *LPopCommand {
	return &LPopCommand{storage}
}

func (c *LPopCommand) Execute(args []resp.Value) resp.Value {
	if len(args) == 0 || len(args) > 2 {
		return core.WrongNumberOfArgumentsError("llen")
	}

	key := args[0].String()
	count := 1
	var err error

	if len(args) == 2 {
		count, err = strconv.Atoi(args[1].String())
		if err != nil {
			return core.InvalidArgumentTypeError()
		}
	}

	entry, exists := c.storage.Get(key)
	if !exists {
		return resp.NewNullBulkString()
	}

	list, isList := entry.(*store.List)
	if !isList {
		return core.WrongTypeOperationError()
	}
	if list.IsEmpty() {
		return resp.NewNullBulkString()
	}

	if count == 1 {
		poppedElement := list.Pop()
		return poppedElement
	}

	i := 1
	poppedElements := []resp.Value{}
	for i <= count && !list.IsEmpty() {
		poppedElements = append(poppedElements, list.Pop())
		i++
	}

	return resp.NewArray(poppedElements)
}
