package main

import (
	"fmt"

	"github.com/md-talim/codecrafters-redis-go/internal/commands"
	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type Redis struct {
	storage  store.Storage
	config   *config.Config
	registry *commands.Registry
}

func NewRedis(storage store.Storage, config *config.Config) *Redis {
	registry := commands.NewRegistry(storage, config)

	return &Redis{
		storage:  storage,
		config:   config,
		registry: registry,
	}
}

func (r *Redis) Evaluate(command resp.Value) resp.Value {
	array, ok := command.(*resp.Array)
	if !ok {
		return resp.NewSimpleError("ERR command must be an array")
	}

	return r.evaluateArray(array)
}

func (r *Redis) evaluateArray(array *resp.Array) resp.Value {
	items := array.Items()
	if len(items) == 0 {
		return nil
	}

	commandName := items[0].String()

	command, exists := r.registry.GetCommand(commandName)
	if !exists {
		return resp.NewSimpleError(fmt.Sprintf("ERR unknown command '%s'", commandName))
	}

	args := items[1:]
	return command.Execute(args)
}
