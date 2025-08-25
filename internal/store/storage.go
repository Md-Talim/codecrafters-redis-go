package store

import (
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/rdb"
)

type Storage interface {
	Set(key string, value any) error
	SetWithExpiry(key string, value any, expiry time.Duration) error
	Get(key string) (any, bool)
	Delete(key string) error
	Keys() []string
}

func New(cfg *config.Config) Storage {
	storage := NewInMemory()
	if cfg.Dir != "" && cfg.DBFilename != "" {
		loadRDBData(storage, cfg.Dir, cfg.DBFilename)
	}
	return storage
}

func loadRDBData(storage *InMemory, dir, filename string) error {
	reader, err := rdb.NewReader(dir, filename)
	if err != nil {
		return err
	}
	if reader == nil {
		return nil
	}
	defer reader.Close()

	data, err := reader.ReadRDB()
	if err != nil {
		return err
	}

	now := time.Now()
	for key, value := range data.Keys {
		// Skip expired keys
		if value.ExpiresAt != nil && now.After(*value.ExpiresAt) {
			continue
		}

		if value.ExpiresAt != nil {
			storage.SetWithExpiry(key, value.Value, value.ExpiresAt.Sub(now))
		} else {
			storage.Set(key, value.Value)
		}
	}

	return nil
}
