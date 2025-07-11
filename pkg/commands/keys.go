package commands

import (
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type KeysCommand struct {
	storage storage.Storage
}

func NewKeysCommand(storage storage.Storage) *KeysCommand {
	return &KeysCommand{storage}
}

func (k *KeysCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("keys")
	}

	if args[0].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	pattern := args[0].Bulk

	switch pattern {
	case "*":
		return k.sendAllKeys()
	default:
		return resp.NewSimpleError("ERR Unsupported pattern for the 'keys' command")
	}
}

func (k *KeysCommand) sendAllKeys() *resp.Value {
	keysList := k.storage.Keys()
	response := []resp.Value{}

	for _, key := range keysList {
		response = append(response, *resp.NewBulkString(key))
	}

	return resp.NewArray(response)
}

func (k *KeysCommand) Name() string {
	return "KEYS"
}
