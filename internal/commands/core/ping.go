package core

import "github.com/md-talim/codecrafters-redis-go/internal/resp"

type PingCommand struct{}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (p *PingCommand) Execute(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewSimpleString("PONG")
	}

	return WrongNumberOfArgumentsError("ping")
}

func (p *PingCommand) Name() string {
	return "PING"
}
