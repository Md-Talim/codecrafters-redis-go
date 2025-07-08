package resp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const CRLF string = "\r\n"

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{reader: bufio.NewReader(reader)}
}

func (p *Parser) Parse() (*Value, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case '*':
		return p.parseArray(line)
	case '+':
		return p.parseSimpleString(line)
	case '$':
		return p.parseBulkString(line)
	default:
		return nil, errors.New("unknown RESP type")
	}
}

func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, CRLF), nil
}

func (p *Parser) parseArray(line string) (*Value, error) {
	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	array := make([]Value, count)
	for i := range count {
		val, err := p.Parse()
		if err != nil {
			return nil, err
		}
		array[i] = *val
	}

	return &Value{Type: Array, Array: array}, nil
}

func (p *Parser) parseSimpleString(line string) (*Value, error) {
	return &Value{Type: SimpleString, String: line[1:]}, nil
}

func (p *Parser) parseBulkString(line string) (*Value, error) {
	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	bulk := make([]byte, length)
	_, err = io.ReadFull(p.reader, bulk)
	if err != nil {
		return nil, err
	}

	p.readLine()

	return &Value{Type: BulkString, Bulk: string(bulk)}, nil
}
