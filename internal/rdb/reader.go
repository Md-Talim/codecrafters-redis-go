package rdb

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Reader struct {
	file *os.File
}

func NewReader(dir, filename string) (*Reader, error) {
	path := filepath.Join(dir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open RDB file: %w", err)
	}

	return &Reader{file}, nil
}

func (r *Reader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func (r *Reader) ReadRDB() (*RDBData, error) {
	if r == nil || r.file == nil {
		return NewRDBData(), nil
	}

	data := NewRDBData()

	if err := r.readHeader(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		b, err := r.readByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch b {
		case 0xFA: // Start of a metadata subsection
			if err := r.skipMetadata(); err != nil {
				return nil, err
			}
		case 0xFE: // Start of a database subsection
			if err := r.readDatabase(data); err != nil {
				return nil, err
			}
		case 0xFF:
			return data, nil
		default:
			return nil, fmt.Errorf("unexpected byte: 0x%02X", b)
		}
	}

	return data, nil
}

func (r *Reader) readHeader() error {
	header := make([]byte, 9) // REDIS0011
	_, err := io.ReadFull(r.file, header)
	if err != nil {
		return err
	}

	expectedHeader := "REDIS0011"
	if string(header) != expectedHeader {
		return fmt.Errorf("invalid header: expected %s, got %s", expectedHeader, string(header))
	}

	return nil
}

// skipMetadata skips the metdata key and the value
func (r *Reader) skipMetadata() error {
	if err := r.skipString(); err != nil {
		return err
	}
	return r.skipString()
}

func (r *Reader) skipString() error {
	_, err := r.readString()
	return err
}

func (r *Reader) readDatabase(data *RDBData) error {
	// Read database index (size encoded)
	_, err := r.readSize()
	if err != nil {
		return err
	}

	// Check for hash table size information
	b, err := r.readByte()
	if err != nil {
		return err
	}

	if b == 0xFB {
		// Skip hash table sizes
		if _, err := r.readSize(); err != nil { // key-value hash table size
			return err
		}
		if _, err := r.readSize(); err != nil { // key-expires hash table size
			return err
		}
	} else {
		r.file.Seek(-1, io.SeekCurrent)
	}

	// Read key-value pairs
	for {
		b, err := r.readByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if b == 0xFE || b == 0xFF || b == 0xFA {
			r.file.Seek(-1, io.SeekCurrent)
			break
		}

		var expiresAt *time.Time

		switch b {
		case 0xFD: // Expire in seconds
			timestamp, err := r.readUint32()
			if err != nil {
				return err
			}
			t := time.Unix(int64(timestamp), 0)
			expiresAt = &t

			b, err = r.readByte()
			if err != nil {
				return err
			}
		case 0xFC: // Expires in milliseconds
			timestamp, err := r.readUint64()
			if err != nil {
				return err
			}
			t := time.Unix(int64(timestamp/1000), int64((timestamp%1000)*1000000))
			expiresAt = &t

			b, err = r.readByte()
			if err != nil {
				return err
			}
		}

		if b != 0x00 {
			return fmt.Errorf("unsupported value type: 0x%02x", b)
		}

		key, err := r.readString()
		if err != nil {
			return err
		}

		value, err := r.readString()
		if err != nil {
			return err
		}

		data.Keys[key] = &RDBValue{
			Value:     value,
			ExpiresAt: expiresAt,
		}
	}

	return nil
}

func (r *Reader) readUint32() (uint32, error) {
	var b [4]byte
	_, err := io.ReadFull(r.file, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:]), nil
}

func (r *Reader) readUint64() (uint64, error) {
	var b [8]byte
	_, err := io.ReadFull(r.file, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}

func (r *Reader) readSize() (uint64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	switch (b & 0xC0) >> 6 {
	case 0:
		return uint64(b & 0x3F), nil
	case 1:
		next, err := r.readByte()
		if err != nil {
			return 0, err
		}
		return uint64((b&0x3F)<<8) | uint64(next), nil
	case 2:
		var bytes [4]byte
		_, err := io.ReadFull(r.file, bytes[:])
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint32(bytes[:])), nil
	case 3:
		return 0, fmt.Errorf("special string encoding not supported in size context")
	}

	return 0, fmt.Errorf("invalid size encoding")
}

func (r *Reader) readString() (string, error) {
	size, err := r.readStringSize()
	if err != nil {
		return "", err
	}

	if size&0xC0 == 0xC0 {
		switch size {
		case 0xC0: // 8-bit integer
			b, err := r.readByte()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%d", b), nil
		case 0xC1: // 16-bit integer
			var bytes [2]byte
			_, err := io.ReadFull(r.file, bytes[:])
			if err != nil {
				return "", nil
			}
			val := binary.LittleEndian.Uint16(bytes[:])
			return fmt.Sprintf("%d", val), nil
		case 0xC2: // 32-bit integer
			var bytes [4]byte
			_, err := io.ReadFull(r.file, bytes[:])
			if err != nil {
				return "", nil
			}
			val := binary.LittleEndian.Uint32(bytes[:])
			return fmt.Sprintf("%d", val), nil
		default:
			return "", fmt.Errorf("unsupported string encoding: 0x%02x", size)
		}
	}

	if size == 0 {
		return "", nil
	}

	bytes := make([]byte, size)
	_, err = io.ReadFull(r.file, bytes)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (r *Reader) readStringSize() (uint64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	switch (b & 0xC0) >> 6 {
	case 0:
		return uint64(b & 0x3F), nil
	case 1:
		next, err := r.readByte()
		if err != nil {
			return 0, err
		}
		return uint64((b&0x3F)<<8) | uint64(next), nil
	case 2:
		var bytes [4]byte
		_, err := io.ReadFull(r.file, bytes[:])
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint32(bytes[:])), nil
	case 3:
		return uint64(b), nil
	}

	return 0, fmt.Errorf("invalid string size encoding")
}

func (r *Reader) readByte() (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r.file, b[:])
	return b[0], err
}
