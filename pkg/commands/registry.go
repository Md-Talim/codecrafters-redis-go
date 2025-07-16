package commands

import (
	"strings"
	"sync"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/interfaces"
	"github.com/md-talim/codecrafters-redis-go/pkg/replication"
)

type Registry struct {
	storage           storage.Storage
	replicaInfo       *replication.Info
	config            *config.Config
	commands          map[string]interfaces.CommandExecutor
	writeCommands     map[string]bool
	handshakeCommands map[string]bool
	mu                sync.RWMutex
}

func NewRegistry(storage storage.Storage, replicaInfo *replication.Info, config *config.Config) *Registry {
	registry := &Registry{
		storage:           storage,
		replicaInfo:       replicaInfo,
		config:            config,
		commands:          make(map[string]interfaces.CommandExecutor),
		writeCommands:     make(map[string]bool),
		handshakeCommands: make(map[string]bool),
	}

	registry.registerCommands()
	registry.registerCommandTypes()

	return registry
}

func (r *Registry) registerCommands() {
	r.commands = map[string]interfaces.CommandExecutor{
		"CONFIG":   NewConfigCommand(r.config),
		"ECHO":     NewEchoCommand(),
		"GET":      NewGetCommand(r.storage),
		"INFO":     NewInfoCommand(r.replicaInfo),
		"KEYS":     NewKeysCommand(r.storage),
		"PING":     NewPingCommand(),
		"PSYNC":    NewPsyncCommand(r.replicaInfo),
		"REPLCONF": NewReplConfCommand(),
		"SET":      NewSetCommand(r.storage),
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
