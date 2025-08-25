package core

import (
	"testing"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/resp"
)

func TestConfigGetDir(t *testing.T) {
	cfg := &config.Config{
		Dir:        "/tmp/redis-files",
		DBFilename: "dump.rdb",
		Port:       "6379",
	}

	cmd := NewConfigCommand(cfg)
	args := []resp.Value{
		resp.NewBulkString("GET"),
		resp.NewBulkString("dir"),
	}

	result := cmd.Execute(args)

	array, isArray := result.(*resp.Array)
	if !isArray {
		t.Errorf("Expected Array")
	}

	items := array.Items()
	if len(items) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(items))
	}

	if items[0].String() != "dir" {
		t.Errorf("Expected 'dir', got %q", items[0].String())
	}

	if items[1].String() != "/tmp/redis-files" {
		t.Errorf("Expected '/tmp/redis-files', got %q", items[1].String())
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
		resp.NewBulkString("GET"),
		resp.NewBulkString("dbfilename"),
	}

	result := cmd.Execute(args)

	array, isArray := result.(*resp.Array)
	if !isArray {
		t.Errorf("Expected Array")
	}

	items := array.Items()
	if len(items) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(items))
	}

	if items[0].String() != "dbfilename" {
		t.Errorf("Expected 'dbfilename', got %q", items[0].String())
	}

	if items[1].String() != "dump.rdb" {
		t.Errorf("Expected 'dump.rdb', got %q", items[1].String())
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
		resp.NewBulkString("GET"),
		resp.NewBulkString("unknown"),
	}

	result := cmd.Execute(args)

	array, isArray := result.(*resp.Array)
	if !isArray {
		t.Errorf("Expected Array")
	}

	if len(array.Items()) != 0 {
		t.Errorf("Expected empty array, got %d elements", len(array.Items()))
	}
}
