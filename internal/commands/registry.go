package commands

import (
	"strings"
	"sync"

	"github.com/md-talim/codecrafters-redis-go/internal/commands/core"
	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type CommandHandler interface {
	Execute([]resp.Value) resp.Value
}

type Registry struct {
	storage  store.Storage
	config   *config.Config
	commands map[string]CommandHandler
	mu       sync.RWMutex
}

func NewRegistry(storage store.Storage, config *config.Config) *Registry {
	registry := &Registry{
		storage:  storage,
		config:   config,
		commands: make(map[string]CommandHandler),
	}

	registry.registerCommands()

	return registry
}

func (r *Registry) registerCommands() {
	r.commands = map[string]CommandHandler{
		"CONFIG": core.NewConfigCommand(r.config),
		"ECHO":   core.NewEchoCommand(),
		"GET":    core.NewGetCommand(r.storage),
		"KEYS":   core.NewKeysCommand(r.storage),
		"PING":   core.NewPingCommand(),
		"SET":    core.NewSetCommand(r.storage),
	}
}

func (r *Registry) GetCommand(name string) (CommandHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, exists := r.commands[strings.ToUpper(name)]
	return cmd, exists
}
