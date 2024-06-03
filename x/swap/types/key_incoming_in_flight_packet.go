package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// IncomingInFlightPacketKeyPrefix is the prefix to retrieve all IncomingInFlightPacket
	IncomingInFlightPacketKeyPrefix = "IncomingInFlightPacket/value/"
)

// IncomingInFlightPacketKey returns the store key to retrieve a IncomingInFlightPacket from the index fields
func IncomingInFlightPacketKey(
	index PacketIndex,
) []byte {
	// var key []byte

	//   indexBytes := []byte(index)
	//   key = append(key, indexBytes...)
	//   key = append(key, []byte("/")...)

	// return key

	bz, err := index.Marshal()
	if err != nil {
		panic(err)
	}
	return append([]byte(IncomingInFlightPacketKeyPrefix), bz...)
}
