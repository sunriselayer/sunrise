package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // TwapKeyPrefix is the prefix to retrieve all Twap
	TwapKeyPrefix = "Twap/value/"
)

// TwapKey returns the store key to retrieve a Twap from the index fields
func TwapKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}