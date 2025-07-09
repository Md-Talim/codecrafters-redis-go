package commands

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func TestConfigGetDir(t *testing.T) {
	cfg := &config.Config{
		Dir:        "/tmp/redis-files",
		DBFilename: "dump.rdb",
		Port:       "6379",
	}

	cmd := NewConfigCommand(cfg)
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "GET"},
		{Type: resp.BulkString, Bulk: "dir"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.Array {
		t.Errorf("Expected Array, got %v", result.Type)
	}

	if len(result.Array) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result.Array))
	}

	if result.Array[0].Bulk != "dir" {
		t.Errorf("Expected 'dir', got %q", result.Array[0].Bulk)
	}

	if result.Array[1].Bulk != "/tmp/redis-files" {
		t.Errorf("Expected '/tmp/redis-files', got %q", result.Array[1].Bulk)
	}
}

func TestConfigGetDBFilename(t *testing.T) {
	cfg := &config.Config{
		Dir:        "/tmp/redis-files",
		DBFilename: "dump.rdb",
		Port:       "6379",
	}

	cmd := NewConfigCommand(cfg)
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "GET"},
		{Type: resp.BulkString, Bulk: "dbfilename"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.Array {
		t.Errorf("Expected Array, got %v", result.Type)
	}

	if len(result.Array) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result.Array))
	}

	if result.Array[0].Bulk != "dbfilename" {
		t.Errorf("Expected 'dbfilename', got %q", result.Array[0].Bulk)
	}

	if result.Array[1].Bulk != "dump.rdb" {
		t.Errorf("Expected 'dump.rdb', got %q", result.Array[1].Bulk)
	}
}

func TestConfigGetUnknownParameter(t *testing.T) {
	cfg := &config.Config{
		Dir:        "/tmp/redis-files",
		DBFilename: "dump.rdb",
		Port:       "6379",
	}

	cmd := NewConfigCommand(cfg)
	args := []resp.Value{
		{Type: resp.BulkString, Bulk: "GET"},
		{Type: resp.BulkString, Bulk: "unknown"},
	}

	result := cmd.Execute(args)

	if result.Type != resp.Array {
		t.Errorf("Expected Array, got %v", result.Type)
	}

	if len(result.Array) != 0 {
		t.Errorf("Expected empty array, got %d elements", len(result.Array))
	}
}
