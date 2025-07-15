package main

import (
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-redis-go/app/server"
	"github.com/md-talim/codecrafters-redis-go/internal/config"
)

func main() {
	config := config.Load()
	redisServer, err := server.NewBuilder(config).Build()
	if err != nil {
		fmt.Printf("Failed to build server: %v\n", err)
		os.Exit(1)
	}

	if err := redisServer.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
