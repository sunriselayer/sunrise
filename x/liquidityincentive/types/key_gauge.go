package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// GaugeKeyPrefix is the prefix to retrieve all Gauge
	GaugeKeyPrefix = "Gauge/value/"
)

// GaugeKey returns the store key to retrieve a Gauge from the index fields
func GaugeKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
