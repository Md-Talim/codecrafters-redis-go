package server

import (
	"github.com/md-talim/codecrafters-redis-go/internal/config"
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

func (s *RedisServer) Start() error {
	if s.config.IsReplica() {
		go s.replicationManager.ConnectToMaster()
	}

	return s.networkListener.Listen(s.connectionHandler.Handle)
}
