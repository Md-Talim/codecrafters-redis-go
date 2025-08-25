package main

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/md-talim/codecrafters-redis-go/internal/resp"
)

type Client struct {
	id    int
	conn  net.Conn
	redis *Redis
}

var clientIDCounter int32

func NewClient(conn net.Conn, redis *Redis) *Client {
	id := int(atomic.AddInt32(&clientIDCounter, 1))
	return &Client{
		id:    id,
		conn:  conn,
		redis: redis,
	}
}

func (c *Client) Handle() {
	defer c.conn.Close()
	fmt.Printf("%d: connected\n", c.id)

	parser := resp.NewParser(c.conn)

	for {
		request, err := parser.Parse()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("%d: parse error: %v\n", c.id, err)
			}
			break
		}

		response := c.redis.Evaluate(request)
		if response == nil {
			fmt.Printf("%d: no response\n", c.id)
			continue
		}

		_, err = c.conn.Write(response.Serialize())
		if err != nil {
			fmt.Printf("%d: write error: %v\n", c.id, err)
			break
		}
	}

	fmt.Printf("%d: disconnected\n", c.id)
}

func (c *Client) ID() int {
	return c.id
}
