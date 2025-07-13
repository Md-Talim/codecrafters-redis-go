package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

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
	replicas    map[string]*ReplicaConnection
	replicasMu  sync.RWMutex
}

type ReplicaConnection struct {
	conn net.Conn
	id   string
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
		replicas:    make(map[string]*ReplicaConnection),
		replicasMu:  sync.RWMutex{},
	}
	server.commands = map[string]Command{
		"CONFIG":   commands.NewConfigCommand(server.config),
		"ECHO":     &commands.EchoCommand{},
		"GET":      commands.NewGetCommand(server.storage),
		"INFO":     commands.NewInfoCommand(server.replicaInfo),
		"KEYS":     commands.NewKeysCommand(server.storage),
		"PING":     &commands.PingCommand{},
		"PSYNC":    commands.NewPsyncCommand(replInfo),
		"REPLCONF": &commands.ReplConfCommand{},
		"SET":      commands.NewSetCommand(server.storage),
	}

	return server
}

func (s *RedisServer) Start() error {
	if s.config.IsReplica() {
		go s.connectToMaster()
	}

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

func (s *RedisServer) connectToMaster() {
	masterHost, masterPort := s.config.GetMasterHostPort()
	if masterHost == "" || masterPort == "" {
		fmt.Println("Invalid master host/port configuration")
		return
	}

	for {
		if err := s.performHandshake(masterHost, masterPort); err != nil {
			fmt.Printf("Failed to connect to master %s:%s. Retrying in 5s...\n", masterHost, masterPort)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}

func (s *RedisServer) performHandshake(host, port string) error {
	masterAddr := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to master: %w", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to master at %s\n", masterAddr)

	// Send PING command
	pingCommand := "*1\r\n$4\r\nPING\r\n"
	_, err = conn.Write([]byte(pingCommand))
	if err != nil {
		return fmt.Errorf("failed to send PING to master: %w", err)
	}
	fmt.Println("Sent PING to master")

	// Read response from master
	parser := resp.NewParser(conn)
	response, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to read PING response: %w", err)
	}

	if response.Type == resp.SimpleString && response.String == "PONG" {
		fmt.Println("Received PONG from master")
	} else {
		fmt.Printf("Unexpected response from master: %+v\n", response)
	}

	// Send 2 REPLCONF command
	replconfListeningPortCommand := "*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n6380\r\n"
	_, err = conn.Write([]byte(replconfListeningPortCommand))
	if err != nil {
		return fmt.Errorf("failed to send REPLCONF to master: %w", err)
	}
	response, err = parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to read REPLCONF response: %w", err)
	}

	if response.Type == resp.SimpleString && response.String == "OK" {
		fmt.Println("Received OK from master")
	} else {
		fmt.Printf("Unexpected response from master: %+v\n", response)
	}

	replconfCapabilitiesCommand := "*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"
	_, err = conn.Write([]byte(replconfCapabilitiesCommand))
	if err != nil {
		return fmt.Errorf("failed to send REPLCONF to master: %w", err)
	}
	response, err = parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to read REPLCONF response: %w", err)
	}

	if response.Type == resp.SimpleString && response.String == "OK" {
		fmt.Println("Received OK from master")
	} else {
		fmt.Printf("Unexpected response from master: %+v\n", response)
	}

	// Send the PSYNC command
	psyncCommand := "*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"
	_, err = conn.Write([]byte(psyncCommand))
	if err != nil {
		return fmt.Errorf("failed to send PSYNC to master: %w", err)
	}

	return nil
}

func (s *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	parser := resp.NewParser(conn)

	var isReplica bool
	var replicaID string

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

		if commandName == "PSYNC" && !isReplica {
			isReplica = true
			replicaID = fmt.Sprintf("replica_%s_%d", conn.RemoteAddr().String(), time.Now().UnixNano())
			s.addReplica(replicaID, conn)
		}

		if cmd, exists := s.commands[commandName]; exists {
			response := cmd.Execute(args)

			if !isReplica || s.isHandshakeCommand(commandName) {
				conn.Write([]byte(response.Serialize()))
			}

			if !isReplica && commandName == "SET" {
				commandArray := make([]resp.Value, len(value.Array))
				copy(commandArray, value.Array)
				s.propagateCommand(commandArray)
				continue
			}
		} else {
			if !isReplica {
				conn.Write([]byte(commands.UnknownCommandError(commandName).Serialize()))
			}
		}
	}

	if isReplica && replicaID != "" {
		s.removeReplica(replicaID)
	}
}

func (s *RedisServer) addReplica(id string, conn net.Conn) {
	s.replicasMu.Lock()
	defer s.replicasMu.Unlock()
	s.replicas[id] = &ReplicaConnection{conn, id}
	fmt.Printf("Replica added: %s (total: %d)\n", id, len(s.replicas))
}

func (s *RedisServer) removeReplica(id string) {
	s.replicasMu.Lock()
	defer s.replicasMu.Unlock()
	delete(s.replicas, id)
	fmt.Printf("Replica removed: %s (total: %d)\n", id, len(s.replicas))
}

func (s *RedisServer) propagateCommand(command []resp.Value) {
	if len(s.replicas) == 0 {
		return
	}

	respCommand := &resp.Value{
		Type:  resp.Array,
		Array: command,
	}
	serializedCommand := respCommand.Serialize()

	for id, replica := range s.replicas {
		_, err := replica.conn.Write([]byte(serializedCommand))
		if err != nil {
			fmt.Printf("Failed to propogate command to replica %s: %v\n", id, err)
		} else {
			fmt.Printf("Propogated command to replica %s: %s\n", id, strings.TrimSpace(serializedCommand))
		}
	}
}

func (s *RedisServer) isHandshakeCommand(command string) bool {
	handshakeCommands := map[string]bool{
		"PING":     true,
		"REPLCONF": true,
		"PSYNC":    true,
	}
	return handshakeCommands[command]
}

type Command interface {
	Execute(args []resp.Value) *resp.Value
	Name() string
}
