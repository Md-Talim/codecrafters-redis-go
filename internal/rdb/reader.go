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
		opCode, err := r.readByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch opCode {
		case OpAux: // Start of a metadata subsection
			if err := r.skipMetadata(); err != nil {
				return nil, err
			}
		case OpSelectDB: // Start of a database subsection
			if err := r.readDatabase(data); err != nil {
				return nil, err
			}
		case OpEOF:
			return data, nil
		default:
			return nil, fmt.Errorf("unexpected byte: 0x%02X", opCode)
		}
	}

	return data, nil
}

func (r *Reader) readHeader() error {
	header := make([]byte, len(RDBHeader))
	_, err := io.ReadFull(r.file, header)
	if err != nil {
		return err
	}

	if string(header) != RDBHeader {
		return fmt.Errorf("invalid header: expected %s, got %s", RDBHeader, string(header))
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
	nextByte, err := r.readByte()
	if err != nil {
		return err
	}

	if nextByte == OpResizeDB {
		// Skip hash table sizes
		if _, err := r.readSize(); err != nil { // key-value hash table size
			return err
		}
		if _, err := r.readSize(); err != nil { // key-expires hash table size
			return err
		}
	} else {
		// Put the byte back since it's not a resize op
		r.file.Seek(-1, io.SeekCurrent)
	}

	// Read key-value pairs
	for {
		opCode, err := r.readByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if opCode == OpSelectDB || opCode == OpEOF || opCode == OpAux {
			r.file.Seek(-1, io.SeekCurrent)
			break
		}

		var expiresAt *time.Time

		switch opCode {
		case OpExpireTime:
			timestamp, err := r.readUint32()
			if err != nil {
				return err
			}
			t := time.Unix(int64(timestamp), 0)
			expiresAt = &t

			opCode, err = r.readByte()
			if err != nil {
				return err
			}
		case OpExpireTimeMS:
			timestamp, err := r.readUint64()
			if err != nil {
				return err
			}
			t := time.Unix(int64(timestamp/1000), int64((timestamp%1000)*1000000))
			expiresAt = &t

			opCode, err = r.readByte()
			if err != nil {
				return err
			}
		}

		// Validate value type
		if opCode != ValueTypeString {
			return fmt.Errorf("unsupported value type: 0x%02x", opCode)
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

func (r *Reader) readSize() (uint64, error) {
	b, err := r.readByte()
	if err != nil {
		return 0, err
	}

	encodingType := (b & SizeEncodingMask) >> 6
	switch encodingType {
	case 0:
		return uint64(b & SizeValueMask), nil
	case 1:
		next, err := r.readByte()
		if err != nil {
			return 0, err
		}
		return (uint64(b&SizeValueMask) << 8) | uint64(next), nil
	case 2:
		val, err := r.readUint32BigEndian()
		if err != nil {
			return 0, err
		}
		return uint64(val), nil
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

	if size&SizeEncodingMask == SizeEncodingMask {
		switch size {
		case StringEnc8BitInt:
			b, err := r.readByte()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%d", b), nil
		case StringEnc16BitInt:
			val, err := r.readUint16()
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%d", val), nil
		case StringEnc32BitInt:
			val, err := r.readUint32()
			if err != nil {
				return "", err
			}
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

	encodingType := (b & SizeEncodingMask) >> 6
	switch encodingType {
	case 0:
		return uint64(b & SizeValueMask), nil
	case 1:
		next, err := r.readByte()
		if err != nil {
			return 0, err
		}
		return (uint64(b&SizeValueMask) << 8) | uint64(next), nil
	case 2:
		val, err := r.readUint32BigEndian()
		if err != nil {
			return 0, err
		}
		return uint64(val), nil
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

func (r *Reader) readUint16() (uint16, error) {
	var bytes [2]byte
	_, err := io.ReadFull(r.file, bytes[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(bytes[:]), nil
}

func (r *Reader) readUint32() (uint32, error) {
	var bytes [4]byte
	_, err := io.ReadFull(r.file, bytes[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(bytes[:]), nil
}

func (r *Reader) readUint32BigEndian() (uint32, error) {
	var bytes [4]byte
	_, err := io.ReadFull(r.file, bytes[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(bytes[:]), nil
}

func (r *Reader) readUint64() (uint64, error) {
	var b [8]byte
	_, err := io.ReadFull(r.file, b[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b[:]), nil
}
