package storage

import (
	"time"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
)

type Storage interface {
	Set(key, value string) error
	SetWithExpiry(key, value string, expiry time.Duration) error
	Get(key string) (string, bool)
	Delete(key string) error
	Keys() []string
}

func NewStorage(cfg *config.Config) Storage {
	if cfg.Dir != "" && cfg.DBFilename != "" {
		return NewRDBStorage(cfg.Dir, cfg.DBFilename)
	}
	return NewMemoryStorage()
}
