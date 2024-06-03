package types

import ()

func NewPacketIndex(
	portId string,
	channelId string,
	sequence uint64,
) PacketIndex {
	return PacketIndex{
		PortId:    portId,
		ChannelId: channelId,
		Sequence:  sequence,
	}
}
