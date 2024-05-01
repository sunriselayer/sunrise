package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// TwapKeyPrefix is the prefix to retrieve all Twap
	TwapKeyPrefix = "Twap/value/"
)

// TwapKey returns the store key to retrieve a Twap from the index fields
func TwapKey(
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
