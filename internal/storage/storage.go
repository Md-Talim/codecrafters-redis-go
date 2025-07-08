package storage

import "time"

type Storage interface {
	Set(key, value string) error
	SetWithExpiry(key, value string, expiry time.Duration) error
	Get(key string) (string, bool)
	Delete(key string) error
}
