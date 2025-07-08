package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/commands"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	memoryStorage := storage.NewMemoryStorage()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn, memoryStorage)
	}
}

func handleConnection(conn net.Conn, storage storage.Storage) {
	defer conn.Close()

	parser := resp.NewParser(conn)

	commandMap := map[string]Command{
		"PING": &commands.PingCommand{},
		"ECHO": &commands.EchoCommand{},
		"SET":  commands.NewSetCommand(storage),
		"GET":  commands.NewGetCommand(storage),
	}

	for {
		value, err := parser.Parse()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Error parsing: %v\n", err)
			}
			break
		}

		if value.Type != resp.Array || len(value.Array) == 0 {
			conn.Write([]byte(resp.NewSimpleError("ERR invalid command format").Serialize()))
			continue
		}

		commandName := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		if cmd, exists := commandMap[commandName]; exists {
			response := cmd.Execute(args)
			conn.Write([]byte(response.Serialize()))
		} else {
			conn.Write([]byte(resp.NewSimpleError("ERR unknown command").Serialize()))
		}
	}
}

type Command interface {
	Execute(args []resp.Value) *resp.Value
	Name() string
}
