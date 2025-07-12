package main

import (
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
)

func main() {
	config := config.Load()
	server := NewRedisServer(config)

	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
