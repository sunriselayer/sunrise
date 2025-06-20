package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetOutgoingInFlightPacket set a specific outgoingInFlightPacket in the store from its index
func (k Keeper) SetOutgoingInFlightPacket(ctx context.Context, outgoingInFlightPacket types.OutgoingInFlightPacket) error {
	err := k.OutgoingInFlightPackets.Set(
		ctx,
		types.OutgoingInFlightPacketKey(outgoingInFlightPacket.Index),
		outgoingInFlightPacket,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetOutgoingInFlightPacket returns a outgoingInFlightPacket from its index
func (k Keeper) GetOutgoingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) (val types.OutgoingInFlightPacket, found bool, err error) {
	key := types.OutgoingInFlightPacketKey(types.NewPacketIndex(srcPortId, srcChannelId, sequence))
	has, err := k.OutgoingInFlightPackets.Has(
		ctx,
		key,
	)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.OutgoingInFlightPackets.Get(ctx, key)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveOutgoingInFlightPacket removes a outgoingInFlightPacket from the store
func (k Keeper) RemoveOutgoingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) error {
	err := k.OutgoingInFlightPackets.Remove(
		ctx,
		types.OutgoingInFlightPacketKey(types.NewPacketIndex(srcPortId, srcChannelId, sequence)),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetOutgoingInFlightPackets returns all outgoingInFlightPacket
func (k Keeper) GetOutgoingInFlightPackets(ctx context.Context) (list []types.OutgoingInFlightPacket, err error) {
	err = k.OutgoingInFlightPackets.Walk(
		ctx,
		nil,
		func(key collections.Triple[string, string, uint64], value types.OutgoingInFlightPacket) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
