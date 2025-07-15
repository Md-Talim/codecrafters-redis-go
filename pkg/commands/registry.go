package commands

import (
	"strings"
	"sync"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/replica"
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/interfaces"
)

type Registry struct {
	commands          map[string]interfaces.CommandExecutor
	writeCommands     map[string]bool
	handshakeCommands map[string]bool
	mu                sync.RWMutex
}

func NewRegistry(storage storage.Storage, replicaInfo *replica.Info, config *config.Config) *Registry {
	registry := &Registry{
		commands:          make(map[string]interfaces.CommandExecutor),
		writeCommands:     make(map[string]bool),
		handshakeCommands: make(map[string]bool),
	}

	registry.registerCommands(storage, replicaInfo, config)
	registry.registerCommandTypes()

	return registry
}

func (r *Registry) registerCommands(storage storage.Storage, replicaInfo *replica.Info, config *config.Config) {
	r.commands = map[string]interfaces.CommandExecutor{
		"CONFIG":   NewConfigCommand(config),
		"ECHO":     NewEchoCommand(),
		"GET":      NewGetCommand(storage),
		"INFO":     NewInfoCommand(replicaInfo),
		"KEYS":     NewKeysCommand(storage),
		"PING":     NewPingCommand(),
		"PSYNC":    NewPsyncCommand(replicaInfo),
		"REPLCONF": NewReplConfCommand(),
		"SET":      NewSetCommand(storage),
	}
}

func (r *Registry) registerCommandTypes() {
	// Write commands
	writeCommands := []string{"SET", "DEL", "INCR", "DECR", "LPUSH", "RPUSH"}
	for _, cmd := range writeCommands {
		r.writeCommands[cmd] = true
	}

	// Handshake commands
	handshakeCommands := []string{"PING", "REPLCONF", "PSYNC"}
	for _, cmd := range handshakeCommands {
		r.handshakeCommands[cmd] = true
	}
}

func (r *Registry) GetCommand(name string) (interfaces.CommandExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, exists := r.commands[strings.ToUpper(name)]
	return cmd, exists
}

func (r *Registry) IsWriteCommand(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.writeCommands[strings.ToUpper(name)]
}

func (r *Registry) IsHandshakeCommand(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.handshakeCommands[strings.ToUpper(name)]
}
