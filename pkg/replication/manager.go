package replication

import (
	"fmt"
	"net"
	"sync"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/pkg/interfaces"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type Manager struct {
	replicaInfo    *Info
	config         *config.Config
	replicas       map[string]*Connection
	mu             sync.RWMutex
	masterConn     net.Conn
	commandHandler interfaces.CommandHandler
}

type Connection struct {
	conn net.Conn
	id   string
}

func NewManager(config *config.Config, handler interfaces.CommandHandler) *Manager {
	return &Manager{
		replicaInfo:    createReplicaInfo(config),
		config:         config,
		replicas:       make(map[string]*Connection),
		commandHandler: handler,
	}
}

func (m *Manager) GetReplicationInfo() *Info {
	return m.replicaInfo
}

func (m *Manager) AddReplica(id string, conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.replicas[id] = &Connection{conn: conn, id: id}
	fmt.Printf("Replica added: %s (total: %d)\n", id, len(m.replicas))
}

func (m *Manager) RemoveReplica(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.replicas, id)
	fmt.Printf("Replica removed: %s (total: %d)\n", id, len(m.replicas))
}

func (m *Manager) PropagateCommand(command []resp.Value) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.replicas) == 0 {
		return
	}

	respCommand := resp.NewArray(command)
	serializedCommand := respCommand.Serialize()

	for id, replica := range m.replicas {
		if err := m.sendToReplica(replica, serializedCommand); err != nil {
			fmt.Printf("Failed to propagate to replica %s: %v\n", id, err)
		}
	}
}

func (m *Manager) sendToReplica(replica *Connection, command string) error {
	_, err := replica.conn.Write([]byte(command))
	return err
}

func (m *Manager) ConnectToMaster() error {
	handshake := NewHandshake(m.config)
	conn, err := handshake.Perform()
	if err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	m.masterConn = conn

	go m.listenForPropagatedCommands()
	return nil
}

func (m *Manager) listenForPropagatedCommands() {
	defer func() {
		if m.masterConn != nil {
			m.masterConn.Close()
			fmt.Println("Master connection closed")
		}
	}()

	parser := resp.NewParser(m.masterConn)
	fmt.Println("Listening for propogated commands from master...")

	for {
		value, err := parser.Parse()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Error parsing propogated command: %v\n", err)
			}
			break
		}

		if err := m.processPropogatedCommand(value); err != nil {
			fmt.Printf("Error processing propogated command: %v\n", err)
		}
	}
}

func (m *Manager) processPropogatedCommand(value *resp.Value) error {
	if value.Type != resp.Array || len(value.Array) == 0 {
		return fmt.Errorf("invalid command format")
	}

	fmt.Printf("Received propogated command: %v\n", value.Array)

	return m.commandHandler.ProcessCommand(value.Array)
}

func createReplicaInfo(config *config.Config) *Info {
	replicaInfo := NewInfo()
	if config.IsReplica() {
		replicaInfo.SetAsSlave()
	}
	return replicaInfo
}
