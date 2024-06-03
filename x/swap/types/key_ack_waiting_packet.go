package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AckWaitingPacketKeyPrefix is the prefix to retrieve all AckWaitingPacket
	AckWaitingPacketKeyPrefix = "AckWaitingPacket/value/"
)

// AckWaitingPacketKey returns the store key to retrieve a AckWaitingPacket from the index fields
func AckWaitingPacketKey(
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
	return append([]byte(AckWaitingPacketKeyPrefix), bz...)
}
