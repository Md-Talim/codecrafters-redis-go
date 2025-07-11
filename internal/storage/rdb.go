package storage

import (
	"sync"
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/rdb"
)

type RDBStorage struct {
	dir      string
	filename string
	isLoaded bool
	memory   *MemoryStorage
	mu       sync.RWMutex
	rdbData  *rdb.RDBData
}

func NewRDBStorage(dir, filename string) *RDBStorage {
	return &RDBStorage{
		dir:      dir,
		filename: filename,
		isLoaded: false,
		memory:   NewMemoryStorage(),
	}
}

func (r *RDBStorage) ensureLoaded() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isLoaded {
		return nil
	}

	reader, err := rdb.NewReader(r.dir, r.filename)
	if err != nil {
		return err
	}
	if reader == nil {
		// File doesn't exist, use empty data
		r.rdbData = rdb.NewRDBData()
		r.isLoaded = true
		return nil
	}
	defer reader.Close()

	data, err := reader.ReadRDB()
	if err != nil {
		return err
	}
	r.rdbData = data

	// Load data into memory storage
	now := time.Now()
	for key, value := range data.Keys {
		// Skip expired keys
		if value.ExpiresAt != nil && now.After(*value.ExpiresAt) {
			continue
		}

		if value.ExpiresAt != nil {
			r.memory.SetWithExpiry(key, value.Value, value.ExpiresAt.Sub(now))
		} else {
			r.memory.Set(key, value.Value)
		}
	}

	r.isLoaded = true
	return nil
}

func (r *RDBStorage) Set(key, value string) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}
	return r.memory.Set(key, value)
}

func (r *RDBStorage) SetWithExpiry(key, value string, expiry time.Duration) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}
	return r.memory.SetWithExpiry(key, value, expiry)
}

func (r *RDBStorage) Get(key string) (string, bool) {
	if err := r.ensureLoaded(); err != nil {
		return "", false
	}
	return r.memory.Get(key)
}

func (r *RDBStorage) Delete(key string) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}
	return r.memory.Delete(key)
}

func (r *RDBStorage) Keys() []string {
	if err := r.ensureLoaded(); err != nil {
		return []string{}
	}
	return r.memory.Keys()
}
