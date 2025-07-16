package server

import (
	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/commands"
	"github.com/md-talim/codecrafters-redis-go/pkg/network"
	"github.com/md-talim/codecrafters-redis-go/pkg/replication"
)

type Builder struct {
	config *config.Config
}

func NewBuilder(config *config.Config) *Builder {
	return &Builder{config: config}
}

func (b *Builder) Build() (*RedisServer, error) {
	storage := storage.New(b.config)

	replicaProcessor := replication.NewReplicaCommandProcessor(nil)
	replicationManager := replication.NewManager(b.config, replicaProcessor)

	commandRegistry := commands.NewRegistry(storage, replicationManager.GetReplicationInfo(), b.config)
	replicaProcessor.SetCommandRegistry(commandRegistry)

	listener := network.NewTCPListener(b.config.Port)
	connectionHandler := network.NewConnectionHandler(commandRegistry, replicationManager)

	return &RedisServer{
		config:             b.config,
		networkListener:    listener,
		connectionHandler:  connectionHandler,
		replicationManager: replicationManager,
		commandRegistry:    commandRegistry,
	}, nil
}
