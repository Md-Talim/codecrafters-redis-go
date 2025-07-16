package interfaces

import (
	"net"

	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type CommandExecutor interface {
	Execute(args []resp.Value) *resp.Value
}

type ReplicationManager interface {
	AddReplica(id string, conn net.Conn)
	RemoveReplica(id string)
	PropagateCommand(command []resp.Value)
	ConnectToMaster() error
}

type ConnectionHandler interface {
	Handle(conn net.Conn)
}

type NetworkListener interface {
	Listen(handler func(net.Conn)) error
	Stop() error
}

type CommandRegistry interface {
	GetCommand(name string) (CommandExecutor, bool)
	IsWriteCommand(name string) bool
	IsHandshakeCommand(name string) bool
}

type CommandHandler interface {
	ProcessCommand(args []resp.Value) error
}
