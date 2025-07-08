package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestServerPing(t *testing.T) {
	go main()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send PING Command
	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	expected := "+PONG\r\n"
	if response != expected {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}
