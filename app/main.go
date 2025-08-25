package main

import (
	"fmt"
	"os"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
)

func main() {
	fmt.Println("Starting Redis server...")

	cfg := config.Load()
	server := NewServer(cfg)

	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
