package rdb

// RDB Op Codes as defined in the RDB file format specification
const (
	OpEOF          = 0xFF // OpEOF marks the end of the RDB file
	OpSelectDB     = 0xFE // OpSelectDB indicates a database selector
	OpExpireTime   = 0xFD // OpExpireTime indicates expire time in seconds
	OpExpireTimeMS = 0xFC // OpExpireTimeMS indicates expire time in milliseconds
	OpResizeDB     = 0xFB // OpResizeDB indicates hash table sizes for main keyspace and expres
	OpAux          = 0xFA // OpAux indicates auxiliary fields (arbitrary key-value settings)
)

// Value type encodings
const (
	ValueTypeString = 0x00
)

// String encoding special values
const (
	StringEnc8BitInt  = 0xC0 // Indicates an 8-bit integer encoding
	StringEnc16BitInt = 0xC1 // Indicates an 16-bit integer encoding
	StringEnc32BitInt = 0xC2 // Indicates an 32-bit integer encoding
	StringEncLZF      = 0xC3 // Indicates LZF compression (not implemented)
)

// Size encoding masks
const (
	SizeEncodingMask = 0xC0 // Used to extract the 2 bits for size encoding type
	SizeValueMask    = 0x3F // Used to extract the remaining 6 bits for size value
)

// RDB file header
const (
	RDBMagicString = "REDIS"
	RDBVersion     = "0011"
	RDBHeader      = RDBMagicString + RDBVersion
)
