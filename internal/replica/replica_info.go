package replica

import (
	"fmt"
	"strings"
	"sync"
)

type Role string

const (
	RoleMaster Role = "master"
	RoleSlave  Role = "slave"
)

type Info struct {
	mu               sync.RWMutex
	role             Role
	masterReplID     string
	masterReplOffset int64
}

func NewInfo() *Info {
	return &Info{
		role:             RoleMaster,
		masterReplID:     generateReplID(),
		masterReplOffset: 0,
	}
}

func (i *Info) SetAsSlave() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.role = RoleSlave
}

func (i *Info) Role() Role {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.role
}

func (i *Info) InfoString() string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	info := []string{
		"role:" + string(i.role),
		"master_replid:" + i.masterReplID,
		"master_repl_offset:" + fmt.Sprintf("%d", i.masterReplOffset),
	}

	return strings.Join(info, "\r\n")
}

func generateReplID() string {
	// Hardcoded repl id
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}
