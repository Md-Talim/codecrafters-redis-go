package main

import (
	"fmt"
	"net"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/store"
)

type Server struct {
	config *config.Config
	redis  *Redis
}

func NewServer(cfg *config.Config) *Server {
	storage := store.New(cfg)
	return &Server{
		config: cfg,
		redis:  NewRedis(storage, cfg),
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("0.0.0.0:%s", s.config.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %w", s.config.Port, err)
	}
	defer listener.Close()

	fmt.Printf("Redis server listening on port %s\n", s.config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		client := NewClient(conn, s.redis)
		go client.Handle()
	}
}
