package rdb

import "time"

type RDBValue struct {
	Value     string
	ExpiresAt *time.Time
}

type RDBData struct {
	Keys map[string]*RDBValue
}

func NewRDBData() *RDBData {
	return &RDBData{
		Keys: make(map[string]*RDBValue),
	}
}
