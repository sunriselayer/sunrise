package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetOutgoingInFlightPacket set a specific outgoingInFlightPacket in the store from its index
func (k Keeper) SetOutgoingInFlightPacket(ctx context.Context, outgoingInFlightPacket types.OutgoingInFlightPacket) {
	err := k.OutgoingInFlightPackets.Set(
		ctx,
		types.OutgoingInFlightPacketsKey(outgoingInFlightPacket.Index),
		outgoingInFlightPacket,
	)
	if err != nil {
		panic(err)
	}
}

// OutgoingInFlightPacket returns a outgoingInFlightPacket from its index
func (k Keeper) GetOutgoingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) (val types.OutgoingInFlightPacket, found bool) {
	key := types.OutgoingInFlightPacketsKey(types.NewPacketIndex(srcPortId, srcChannelId, sequence))
	has, err := k.OutgoingInFlightPackets.Has(
		ctx,
		key,
	)
	if err != nil {
		panic(err)
	}

	if !has {
		return val, false
	}

	val, err = k.OutgoingInFlightPackets.Get(ctx, key)
	if err != nil {
		panic(err)
	}

	return val, true
}

// RemoveOutgoingInFlightPacket removes a outgoingInFlightPacket from the store
func (k Keeper) RemoveOutgoingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) {
	err := k.OutgoingInFlightPackets.Remove(
		ctx,
		types.OutgoingInFlightPacketsKey(types.NewPacketIndex(srcPortId, srcChannelId, sequence)),
	)
	if err != nil {
		panic(err)
	}
}

// OutgoingInFlightPackets returns all outgoingInFlightPacket
func (k Keeper) GetOutgoingInFlightPackets(ctx context.Context) (list []types.OutgoingInFlightPacket) {
	err := k.OutgoingInFlightPackets.Walk(
		ctx,
		nil,
		func(key collections.Triple[string, string, uint64], value types.OutgoingInFlightPacket) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}
