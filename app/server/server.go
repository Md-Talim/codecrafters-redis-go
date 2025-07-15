package server

import (
	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/replica"
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/commands"
	"github.com/md-talim/codecrafters-redis-go/pkg/network"
	"github.com/md-talim/codecrafters-redis-go/pkg/replication"
)

type RedisServer struct {
	config             *config.Config
	networkListener    *network.TCPListener
	connectionHandler  *network.ConnectionHandler
	replicationManager *replication.Manager
	commandRegistry    *commands.Registry
}

func NewRedisServer(config *config.Config) *RedisServer {
	storage := storage.New(config)
	replicaInfo := replica.NewInfo()

	if config.IsReplica() {
		replicaInfo.SetAsSlave()
	}

	commandsRegistry := commands.NewRegistry(storage, replicaInfo, config)
	replicationManager := replication.NewManager(replicaInfo, config)

	return &RedisServer{
		config:             config,
		networkListener:    network.NewTCPListener(config.Port),
		replicationManager: replicationManager,
		commandRegistry:    commandsRegistry,
		connectionHandler:  network.NewConnectionHandler(commandsRegistry, replicationManager),
	}
}

func (s *RedisServer) Start() error {
	if s.config.IsReplica() {
		go s.replicationManager.ConnectToMaster()
	}

	return s.networkListener.Listen(s.connectionHandler.Handle)
}
