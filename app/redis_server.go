package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/replica"
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/commands"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type RedisServer struct {
	storage     storage.Storage
	config      *config.Config
	replicaInfo *replica.Info
	commands    map[string]Command
}

func NewRedisServer(config *config.Config) *RedisServer {
	replInfo := replica.NewInfo()
	if config.IsReplica() {
		replInfo.SetAsSlave()
	}

	dataStorage := storage.New(config)

	server := &RedisServer{
		storage:     dataStorage,
		config:      config,
		replicaInfo: replInfo,
	}
	server.commands = map[string]Command{
		"CONFIG": commands.NewConfigCommand(server.config),
		"ECHO":   &commands.EchoCommand{},
		"GET":    commands.NewGetCommand(server.storage),
		"INFO":   commands.NewInfoCommand(server.replicaInfo),
		"KEYS":   commands.NewKeysCommand(server.storage),
		"PING":   &commands.PingCommand{},
		"SET":    commands.NewSetCommand(server.storage),
	}

	return server
}

func (s *RedisServer) Start() error {
	address := fmt.Sprintf("0.0.0.0:%s", s.config.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Printf("Redis server listening on %s\n", address)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	parser := resp.NewParser(conn)

	for {
		value, err := parser.Parse()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Error parsing: %v\n", err)
			}
			break
		}

		if value.Type != resp.Array || len(value.Array) == 0 {
			conn.Write([]byte(commands.InvalidCommandFormatError().Serialize()))
			continue
		}

		commandName := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		if cmd, exists := s.commands[commandName]; exists {
			response := cmd.Execute(args)
			conn.Write([]byte(response.Serialize()))
		} else {
			conn.Write([]byte(commands.UnknownCommandError(commandName).Serialize()))
		}
	}
}

type Command interface {
	Execute(args []resp.Value) *resp.Value
	Name() string
}
