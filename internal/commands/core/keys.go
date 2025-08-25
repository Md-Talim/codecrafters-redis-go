package core

import (
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type KeysCommand struct {
	storage store.Storage
}

func NewKeysCommand(storage store.Storage) *KeysCommand {
	return &KeysCommand{storage}
}

func (k *KeysCommand) Execute(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("keys")
	}

	pattern := args[0].String()

	switch pattern {
	case "*":
		return k.sendAllKeys()
	default:
		return resp.NewSimpleError("ERR Unsupported pattern for the 'keys' command")
	}
}

func (k *KeysCommand) sendAllKeys() resp.Value {
	keysList := k.storage.Keys()
	response := []resp.Value{}

	for _, key := range keysList {
		response = append(response, resp.NewBulkString(key))
	}

	return resp.NewArray(response)
}

func (k *KeysCommand) Name() string {
	return "KEYS"
}
