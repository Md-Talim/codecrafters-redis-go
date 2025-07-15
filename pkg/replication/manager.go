package replication

import (
	"fmt"
	"net"
	"sync"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/replica"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type Manager struct {
	replicaInfo *replica.Info
	config      *config.Config
	replicas    map[string]*Connection
	mu          sync.RWMutex
}

type Connection struct {
	conn net.Conn
	id   string
}

func NewManager(replicaInfo *replica.Info, config *config.Config) *Manager {
	return &Manager{
		replicaInfo: replicaInfo,
		config:      config,
		replicas:    make(map[string]*Connection),
	}
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
	return handshake.Perform()
}
