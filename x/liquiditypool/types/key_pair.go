package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PairKeyPrefix is the prefix to retrieve all Pair
	PairKeyPrefix = "Pair/value/"
)

// PairKey returns the store key to retrieve a Pair from the index fields
func PairKey(
	baseDenom string,
	quoteDenom string,
) []byte {
	var key []byte

	baseDenomBytes := []byte(baseDenom)
	quoteDenomBytes := []byte(quoteDenom)
	key = append(key, baseDenomBytes...)
	key = append(key, []byte("/")...)
	key = append(key, quoteDenomBytes...)

	return key
}
