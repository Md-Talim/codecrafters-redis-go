package commands

import (
	"strings"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type ConfigCommand struct {
	config *config.Config
}

func NewConfigCommand(config *config.Config) *ConfigCommand {
	return &ConfigCommand{config}
}

func (c *ConfigCommand) Execute(args []resp.Value) *resp.Value {
	if len(args) < 2 {
		return WrongNumberOfArgumentsError("config")
	}

	if args[0].Type != resp.BulkString || args[1].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	subcommand := strings.ToUpper(args[0].Bulk)

	switch subcommand {
	case "GET":
		return c.handleGet(args[1:])
	default:
		return resp.NewSimpleError("ERR Unknown CONFIG subcommand or wrong number of arguments")
	}
}

func (c *ConfigCommand) handleGet(args []resp.Value) *resp.Value {
	if len(args) != 1 {
		return WrongNumberOfArgumentsError("config get")
	}

	if args[0].Type != resp.BulkString {
		return InvalidArgumentTypeError()
	}

	param := strings.ToLower(args[0].Bulk)
	value, exists := c.config.GetParameter(param)
	if !exists {
		return resp.NewArray([]resp.Value{})
	}

	result := []resp.Value{
		*resp.NewBulkString(param),
		*resp.NewBulkString(value),
	}

	return resp.NewArray(result)
}

func (c *ConfigCommand) Name() string {
	return "CONFIG"
}
