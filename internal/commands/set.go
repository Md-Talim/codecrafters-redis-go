package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type SetCommand struct {
	storage store.Storage
}

func NewSetCommand(storage store.Storage) *SetCommand {
	return &SetCommand{storage}
}

func (s *SetCommand) Execute(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return WrongNumberOfArgumentsError("set")
	}

	key := args[0].String()
	value := args[1].String()

	var expiry time.Duration
	var hasExpiry bool

	for i := 2; i < len(args); i++ {
		arg := strings.ToUpper(args[i].String())
		switch arg {
		case "PX":
			if i+1 >= len(args) {
				return SyntaxError()
			}

			milliseconds, err := strconv.ParseInt(args[i+1].String(), 10, 64)
			if err != nil || milliseconds <= 0 {
				return InvalidExpireTimeError("set")
			}

			expiry = time.Duration(milliseconds) * time.Millisecond
			hasExpiry = true
			i++
		default:
			return SyntaxError()
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
