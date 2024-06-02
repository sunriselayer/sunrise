package types

import "encoding/binary"

import "fmt"

var _ binary.ByteOrder

const (
	// InFlightPacketKeyPrefix is the prefix to retrieve all InFlightPacket
	InFlightPacketKeyPrefix = "InFlightPacket/value/"
)

// InFlightPacketKey returns the store key to retrieve a InFlightPacket from the index fields
func InFlightPacketKey(
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) []byte {
	index := fmt.Sprintf("%s/%s/%d", srcPortId, srcChannelId, sequence)

	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
