package commands

import (
	"strconv"
	"strings"
	"time"

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
	if len(args) < 2 {
		return resp.NewSimpleError("ERR wrong number of arguments for 'set' command")
	}

	if args[0].Type != resp.BulkString || args[1].Type != resp.BulkString {
		return resp.NewSimpleError("Err invalid argument type")
	}

	key := args[0].Bulk
	value := args[1].Bulk

	var expiry time.Duration
	var hasExpiry bool

	for i := 2; i < len(args); i++ {
		if args[i].Type != resp.BulkString {
			return resp.NewSimpleError("ERR invalid argument type")
		}

		arg := strings.ToUpper(args[i].Bulk)
		switch arg {
		case "PX":
			if i+1 >= len(args) {
				return resp.NewSimpleError("ERR syntax error")
			}

			if args[i+1].Type != resp.BulkString {
				return resp.NewSimpleError("ERR invalid argument type")
			}

			milliseconds, err := strconv.ParseInt(args[i+1].Bulk, 10, 64)
			if err != nil || milliseconds <= 0 {
				return resp.NewSimpleError("ERR invalid expire time in 'set' command")
			}

			expiry = time.Duration(milliseconds) * time.Millisecond
			hasExpiry = true
			i++
		default:
			return resp.NewSimpleError("ERR syntax error")
		}
	}

	var err error
	if hasExpiry {
		err = s.storage.SetWithExpiry(key, value, expiry)
	} else {
		err = s.storage.Set(key, value)
	}
	if err != nil {
		return resp.NewSimpleError("ERR failed to set key")
	}

	return resp.NewSimpleString("OK")
}

func (s *SetCommand) Name() string {
	return "SET"
}
