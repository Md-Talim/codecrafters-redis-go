package commands

import (
	"strings"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
)

type ConfigCommand struct {
	config *config.Config
}

func NewConfigCommand(config *config.Config) *ConfigCommand {
	return &ConfigCommand{config}
}

func (c *ConfigCommand) Execute(args []resp.Value) resp.Value {
	if len(args) < 2 {
		return WrongNumberOfArgumentsError("config")
	}

	subcommand := strings.ToUpper(args[0].String())

	switch subcommand {
	case "GET":
		return c.handleGet(args[1:])
	default:
		return resp.NewSimpleError("ERR Unknown CONFIG subcommand or wrong number of arguments")
	}
}

func (c *ConfigCommand) handleGet(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("config get")
	}

	param := strings.ToLower(args[0].String())
	value, exists := c.config.GetParameter(param)
	if !exists {
		return resp.NewArray([]resp.Value{})
	}

	result := []resp.Value{
		resp.NewBulkString(param),
		resp.NewBulkString(value),
	}

	return resp.NewArray(result)
}

func (c *ConfigCommand) Name() string {
	return "CONFIG"
}
