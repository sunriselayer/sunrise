package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// OutgoingInFlightPacketKeyPrefix is the prefix to retrieve all OutgoingInFlightPacket
	OutgoingInFlightPacketKeyPrefix = "OutgoingInFlightPacket/value/"
)

// OutgoingInFlightPacketKey returns the store key to retrieve a OutgoingInFlightPacket from the index fields
func OutgoingInFlightPacketKey(
	index PacketIndex,
) []byte {
	// var key []byte

	// indexBytes := []byte(index)
	// key = append(key, indexBytes...)
	// key = append(key, []byte("/")...)

	// return key

	bz, err := index.Marshal()
	if err != nil {
		panic(err)
	}
	return append([]byte(OutgoingInFlightPacketKeyPrefix), bz...)
}
