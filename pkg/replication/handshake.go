package replication

import (
	"fmt"
	"net"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

type HandshakeStep struct {
	Name             string
	Command          string
	ExpectedResponse func(*resp.Value) bool
	OnSuccess        func()
}

type Handshake struct {
	config *config.Config
	steps  []HandshakeStep
}

func NewHandshake(config *config.Config) *Handshake {
	h := &Handshake{config: config}
	h.buildSteps()
	return h
}

func (h *Handshake) buildSteps() {
	h.steps = []HandshakeStep{
		{
			Name:    "PING",
			Command: "*1\r\n$4\r\nPING\r\n",
			ExpectedResponse: func(response *resp.Value) bool {
				return response.Type == resp.SimpleString && response.String == "PONG"
			},
			OnSuccess: func() { fmt.Println("PING successful") },
		},
		{
			Name:    "REPLCONF listening-port",
			Command: h.buildReplConfPortCommand(),
			ExpectedResponse: func(expected *resp.Value) bool {
				return expected.Type == resp.SimpleString && expected.String == "OK"
			},
			OnSuccess: func() { fmt.Println("REPLCONF listening-port successful") },
		},
		{
			Name:    "REPLCONF capa",
			Command: "*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n",
			ExpectedResponse: func(expected *resp.Value) bool {
				return expected.Type == resp.SimpleString && expected.String == "OK"
			},
			OnSuccess: func() { fmt.Println("REPLCONF capa successful") },
		},
		{
			Name:    "PSYNC",
			Command: "*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n",
			ExpectedResponse: func(expected *resp.Value) bool {
				return expected.Type == resp.SimpleString &&
					len(expected.String) > 11 &&
					expected.String[:11] == "FULLRESYNC "
			},
			OnSuccess: func() { fmt.Println("PSYNC successful") },
		},
	}
}

func (h *Handshake) buildReplConfPortCommand() string {
	port := h.config.Port
	return fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$%d\r\n%s\r\n",
		len(port), port)
}

func (h *Handshake) Perform() (net.Conn, error) {
	masterHost, masterPort := h.config.GetMasterHostPort()
	if masterHost == "" || masterPort == "" {
		return nil, fmt.Errorf("invalid master configuration")
	}

	conn, err := h.connectToMaster(masterHost, masterPort)
	if err != nil {
		return nil, err
	}

	parser := resp.NewParser(conn)

	for _, step := range h.steps {
		if err := h.executeStep(conn, parser, step); err != nil {
			return nil, fmt.Errorf("step %s failed: %w", step.Name, err)
		}
	}

	return conn, nil
}

func (h *Handshake) connectToMaster(host, port string) (net.Conn, error) {
	masterAddr := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master %s: %w", masterAddr, err)
	}

	fmt.Printf("Connected to master at %s\n", masterAddr)
	return conn, nil
}

func (h *Handshake) executeStep(conn net.Conn, parser *resp.Parser, step HandshakeStep) error {
	// Send command
	if _, err := conn.Write([]byte(step.Command)); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// Read response
	response, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Validate response
	if !step.ExpectedResponse(response) {
		return fmt.Errorf("unexpected response: %+v", response)
	}

	// Execute success callback
	step.OnSuccess()
	return nil
}
