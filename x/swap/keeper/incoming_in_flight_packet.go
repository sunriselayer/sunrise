package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetIncomingInFlightPacket set a specific incomingPacket in the store from its index
func (k Keeper) SetIncomingInFlightPacket(ctx context.Context, incomingPacket types.IncomingInFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IncomingInFlightPacketKeyPrefix))
	b := k.cdc.MustMarshal(&incomingPacket)
	store.Set(types.IncomingInFlightPacketKey(
		incomingPacket.Index,
	), b)
}

// GetIncomingInFlightPacket returns a incomingPacket from its index
func (k Keeper) GetIncomingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) (val types.IncomingInFlightPacket, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IncomingInFlightPacketKeyPrefix))

	b := store.Get(types.IncomingInFlightPacketKey(
		types.NewPacketIndex(
			srcPortId,
			srcChannelId,
			sequence,
		),
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveIncomingInFlightPacket removes a incomingPacket from the store
func (k Keeper) RemoveIncomingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IncomingInFlightPacketKeyPrefix))
	store.Delete(types.IncomingInFlightPacketKey(
		types.NewPacketIndex(
			srcPortId,
			srcChannelId,
			sequence,
		),
	))
}

// GetIncomingInFlightPackets returns all incomingPacket
func (k Keeper) GetIncomingInFlightPackets(ctx context.Context) (list []types.IncomingInFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IncomingInFlightPacketKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.IncomingInFlightPacket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
