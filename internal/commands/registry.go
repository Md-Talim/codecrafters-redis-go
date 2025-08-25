package commands

import (
	"strings"
	"sync"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type Command interface {
	Name() string
	Execute([]resp.Value) resp.Value
}

type Registry struct {
	storage  store.Storage
	config   *config.Config
	commands map[string]Command
	mu       sync.RWMutex
}

func NewRegistry(storage store.Storage, config *config.Config) *Registry {
	registry := &Registry{
		storage:  storage,
		config:   config,
		commands: make(map[string]Command),
	}

	registry.registerCommands()

	return registry
}

func (r *Registry) registerCommands() {
	r.commands = map[string]Command{
		"CONFIG": NewConfigCommand(r.config),
		"ECHO":   NewEchoCommand(),
		"GET":    NewGetCommand(r.storage),
		"KEYS":   NewKeysCommand(r.storage),
		"PING":   NewPingCommand(),
		"SET":    NewSetCommand(r.storage),
	}
}

func (r *Registry) GetCommand(name string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, exists := r.commands[strings.ToUpper(name)]
	return cmd, exists
}
