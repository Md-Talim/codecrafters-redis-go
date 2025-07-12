package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/md-talim/codecrafters-redis-go/internal/config"
	"github.com/md-talim/codecrafters-redis-go/internal/storage"
	"github.com/md-talim/codecrafters-redis-go/pkg/commands"
	"github.com/md-talim/codecrafters-redis-go/pkg/resp"
)

func main() {
	rdbConfig := config.Load()

	address := fmt.Sprintf("0.0.0.0:%s", rdbConfig.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	dataStorage := storage.New(rdbConfig)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn, dataStorage, rdbConfig)
	}
}

func handleConnection(conn net.Conn, storage storage.Storage, rdbConfig *config.Config) {
	defer conn.Close()

	parser := resp.NewParser(conn)

	commandMap := map[string]Command{
		"PING":   &commands.PingCommand{},
		"ECHO":   &commands.EchoCommand{},
		"SET":    commands.NewSetCommand(storage),
		"GET":    commands.NewGetCommand(storage),
		"KEYS":   commands.NewKeysCommand(storage),
		"INFO":   commands.NewInfoCommand(rdbConfig),
		"CONFIG": commands.NewConfigCommand(rdbConfig),
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
			conn.Write([]byte(commands.InvalidCommandFormatError().Serialize()))
			continue
		}

		commandName := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		if cmd, exists := commandMap[commandName]; exists {
			response := cmd.Execute(args)
			conn.Write([]byte(response.Serialize()))
		} else {
			conn.Write([]byte(commands.UnknownCommandError(commandName).Serialize()))
		}
	}
}

type Command interface {
	Execute(args []resp.Value) *resp.Value
	Name() string
}
