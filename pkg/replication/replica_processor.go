package replication

import (
	"fmt"
	"strings"

	"github.com/md-talim/codecrafters-redis-go/pkg/interfaces"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type ReplicaCommandProcessor struct {
	commandRegistry interfaces.CommandRegistry
}

func NewReplicaCommandProcessor(registry interfaces.CommandRegistry) *ReplicaCommandProcessor {
	return &ReplicaCommandProcessor{
		commandRegistry: registry,
	}
}

func (p *ReplicaCommandProcessor) SetCommandRegistry(registry interfaces.CommandRegistry) {
	p.commandRegistry = registry
}

func (p *ReplicaCommandProcessor) ProcessCommand(command []resp.Value) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command")
	}

	commandName := strings.ToUpper(command[0].Bulk)
	args := command[1:]

	cmd, exists := p.commandRegistry.GetCommand(commandName)
	if !exists {
		return nil
	}

	cmd.Execute(args)

	return nil
}
