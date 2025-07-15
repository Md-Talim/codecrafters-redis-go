package server

import (
	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/replica"
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
	// Create core components
	storage := storage.New(b.config)
	replicaInfo := b.createReplicaInfo()

	// Create managers
	replicationManager := replication.NewManager(replicaInfo, b.config)
	commandRegistry := commands.NewRegistry(storage, replicaInfo, b.config)

	// Create network components
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

func (b *Builder) createReplicaInfo() *replica.Info {
	replicaInfo := replica.NewInfo()
	if b.config.IsReplica() {
		replicaInfo.SetAsSlave()
	}
	return replicaInfo
}
