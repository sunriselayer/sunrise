package types

import ()

func NewPacketIndex(
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) PacketIndex {
	return PacketIndex{
		SrcPortId:    srcPortId,
		SrcChannelId: srcChannelId,
		Sequence:     sequence,
	}
}
