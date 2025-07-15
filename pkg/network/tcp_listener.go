package network

import (
	"fmt"
	"net"
)

type TCPListener struct {
	port     string
	listener net.Listener
}

func NewTCPListener(port string) *TCPListener {
	return &TCPListener{port: port}
}

func (t *TCPListener) Listen(handler func(net.Conn)) error {
	address := fmt.Sprintf("0.0.0.0:%s", t.port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %w", t.port, err)
	}

	t.listener = listener
	fmt.Printf("Redis server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handler(conn)
	}
}

func (t *TCPListener) Stop() error {
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}
