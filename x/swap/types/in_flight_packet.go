package types

import ()

func NewInFlightPacketIndex(
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) InFlightPacketIndex {
	return InFlightPacketIndex{
		SrcPortId:    srcPortId,
		SrcChannelId: srcChannelId,
		Sequence:     sequence,
	}
}
