package replica

import (
	"strings"
	"sync"
)

type Role string

const (
	RoleMaster Role = "master"
	RoleSlave  Role = "slave"
)

type Info struct {
	mu   sync.RWMutex
	role Role
}

func NewInfo() *Info {
	return &Info{
		role: RoleMaster,
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
	}

	return strings.Join(info, "\r\n")
}
