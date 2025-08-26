package list

import (
	"strconv"

	"github.com/md-talim/codecrafters-redis-go/internal/commands/core"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type LRangeCommand struct {
	storage store.Storage
}

func NewLRangeCommand(storage store.Storage) *LRangeCommand {
	return &LRangeCommand{storage}
}

func (c *LRangeCommand) Execute(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return core.WrongNumberOfArgumentsError("lrange")
	}

	key := args[0].String()

	start, err := strconv.Atoi(args[1].String())
	if err != nil {
		return core.ValueNotIntegerError()
	}
	stop, err := strconv.Atoi(args[2].String())
	if err != nil {
		return core.ValueNotIntegerError()
	}

	entry, exists := c.storage.Get(key)
	if !exists {
		return resp.NewArray([]resp.Value{})
	}

	list, isList := entry.(*store.List)
	if !isList {
		return core.WrongTypeOperationError()
	}

	size := list.Size()

	if start < 0 {
		start = max(size+start, 0)
	}
	if stop < 0 {
		stop = size + stop
	}

	if start >= size || start > stop {
		return resp.NewArray([]resp.Value{})
	}

	if stop >= size {
		stop = size - 1
	}

	result := list.Range(start, stop+1)

	return resp.NewArray(result)
}
