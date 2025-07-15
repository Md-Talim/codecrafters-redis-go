package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/md-talim/codecrafters-redis-go/pkg/interfaces"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type ConnectionHandler struct {
	commandRegistry    interfaces.CommandRegistry
	replicationManager interfaces.ReplicationManager
}

func NewConnectionHandler(registry interfaces.CommandRegistry, replManager interfaces.ReplicationManager) *ConnectionHandler {
	return &ConnectionHandler{
		commandRegistry:    registry,
		replicationManager: replManager,
	}
}

func (h *ConnectionHandler) Handle(conn net.Conn) {
	defer conn.Close()

	session := NewConnectionSession(conn, h.commandRegistry, h.replicationManager)
	session.ProcessCommands()
}

// Connection Session - handles individual connection lifecycle
type ConnectionSession struct {
	conn               net.Conn
	parser             *resp.Parser
	commandRegistry    interfaces.CommandRegistry
	replicationManager interfaces.ReplicationManager
	isReplica          bool
	replicaID          string
}

func NewConnectionSession(conn net.Conn, registry interfaces.CommandRegistry, replManager interfaces.ReplicationManager) *ConnectionSession {
	return &ConnectionSession{
		conn:               conn,
		parser:             resp.NewParser(conn),
		commandRegistry:    registry,
		replicationManager: replManager,
	}
}

func (s *ConnectionSession) ProcessCommands() {
	defer s.cleanup()

	for {
		if !s.processNextCommand() {
			break
		}
	}
}

func (s *ConnectionSession) processNextCommand() bool {
	value, err := s.parser.Parse()
	if err != nil {
		if err.Error() != "EOF" {
			fmt.Printf("Error parsing: %v\n", err)
		}
		return false
	}

	if !s.isValidCommand(value) {
		s.sendError("Invalid command format")
		return true
	}

	return s.executeCommand(value)
}

func (s *ConnectionSession) executeCommand(value *resp.Value) bool {
	commandName := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]

	s.handleReplicaIdentification(commandName)

	cmd, exists := s.commandRegistry.GetCommand(commandName)
	if !exists {
		if !s.isReplica {
			s.sendError(fmt.Sprintf("Unknown command: %s", commandName))
		}
		return true
	}

	response := cmd.Execute(args)
	s.handleResponse(commandName, response)
	s.handleCommandPropagation(commandName, value.Array)

	return true
}

func (s *ConnectionSession) handleReplicaIdentification(commandName string) {
	if commandName == "PSYNC" && !s.isReplica {
		s.isReplica = true
		s.replicaID = fmt.Sprintf("replica_%s_%d", s.conn.RemoteAddr().String(), time.Now().UnixNano())
		s.replicationManager.AddReplica(s.replicaID, s.conn)
	}
}

func (s *ConnectionSession) handleResponse(commandName string, response *resp.Value) {
	shouldSendResponse := !s.isReplica || s.commandRegistry.IsHandshakeCommand(commandName)

	if shouldSendResponse {
		s.conn.Write([]byte(response.Serialize()))
	}
}

func (s *ConnectionSession) handleCommandPropagation(commandName string, command []resp.Value) {
	if !s.isReplica && s.commandRegistry.IsWriteCommand(commandName) {
		s.replicationManager.PropagateCommand(command)
	}
}

func (s *ConnectionSession) cleanup() {
	if s.isReplica && s.replicaID != "" {
		s.replicationManager.RemoveReplica(s.replicaID)
	}
}

func (s *ConnectionSession) isValidCommand(value *resp.Value) bool {
	return value.Type == resp.Array && len(value.Array) > 0
}

func (s *ConnectionSession) sendError(message string) {
	errorResponse := resp.NewSimpleError(message)
	s.conn.Write([]byte(errorResponse.Serialize()))
}
