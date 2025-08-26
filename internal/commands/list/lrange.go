package list

import (
	"strconv"

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
		return WrongNumberOfArgumentsError("lrange")
	}

	key := args[0].String()

	// TODO: return error if start and stop are not integers
	start, _ := strconv.Atoi(args[1].String())
	stop, _ := strconv.Atoi(args[2].String())

	entry, exists := c.storage.Get(key)
	if !exists {
		return resp.NewArray([]resp.Value{})
	}

	list, isList := entry.(*store.List)
	if !isList {
		return WrongTypeOperationError()
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
	responseArray := []resp.Value{}

	for _, item := range result {
		responseArray = append(responseArray, resp.NewBulkString(item.(string)))
	}

	return resp.NewArray(responseArray)
}
