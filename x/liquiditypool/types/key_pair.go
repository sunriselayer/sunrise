package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // PairKeyPrefix is the prefix to retrieve all Pair
	PairKeyPrefix = "Pair/value/"
)

// PairKey returns the store key to retrieve a Pair from the index fields
func PairKey(
index string,
) []byte {
	var key []byte
    
    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)
    
	return key
}