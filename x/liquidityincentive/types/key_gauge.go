package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// GaugeKeyPrefix is the prefix to retrieve all Gauge
	GaugeKeyPrefix = "Gauge/value/"
)

func GaugeKeyPrefixByPreviousEpochId(previousEpochId uint64) []byte {
	var key []byte

	previoudEpochIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(previoudEpochIdBytes, previousEpochId)

	key = append(key, []byte(GaugeKeyPrefix)...)
	key = append(key, previoudEpochIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

// GaugeKey returns the store key to retrieve a Gauge from the index fields
func GaugeKey(previousEpochId uint64, poolId uint64) []byte {
	var key []byte

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)

	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
