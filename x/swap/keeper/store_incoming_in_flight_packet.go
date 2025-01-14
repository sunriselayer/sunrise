package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetIncomingInFlightPacket set a specific incomingPacket in the store from its index
func (k Keeper) SetIncomingInFlightPacket(ctx context.Context, incomingPacket types.IncomingInFlightPacket) {
	err := k.IncomingInFlightPackets.Set(
		ctx,
		types.IncomingInFlightPacketKey(incomingPacket.Index),
		incomingPacket,
	)
	if err != nil {
		panic(err)
	}
}

// GetIncomingInFlightPacket returns a incomingPacket from its index
func (k Keeper) GetIncomingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) (val types.IncomingInFlightPacket, found bool) {
	key := types.IncomingInFlightPacketKey(types.NewPacketIndex(srcPortId, srcChannelId, sequence))
	has, err := k.IncomingInFlightPackets.Has(
		ctx,
		key,
	)
	if err != nil {
		panic(err)
	}

	if !has {
		return val, false
	}

	val, err = k.IncomingInFlightPackets.Get(ctx, key)
	if err != nil {
		panic(err)
	}

	return val, true
}

// RemoveIncomingInFlightPacket removes a incomingPacket from the store
func (k Keeper) RemoveIncomingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) {
	err := k.IncomingInFlightPackets.Remove(
		ctx,
		types.IncomingInFlightPacketKey(types.NewPacketIndex(srcPortId, srcChannelId, sequence)),
	)
	if err != nil {
		panic(err)
	}
}

// GetIncomingInFlightPackets returns all incomingPacket
func (k Keeper) GetIncomingInFlightPackets(ctx context.Context) (list []types.IncomingInFlightPacket) {
	err := k.IncomingInFlightPackets.Walk(
		ctx,
		nil,
		func(key collections.Triple[string, string, uint64], value types.IncomingInFlightPacket) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}
