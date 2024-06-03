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

func (index PacketIndex) Equal(other PacketIndex) bool {
	return index.PortId == other.PortId &&
		index.ChannelId == other.ChannelId &&
		index.Sequence == other.Sequence
}
