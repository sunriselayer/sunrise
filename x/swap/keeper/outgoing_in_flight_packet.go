package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetOutgoingInFlightPacket set a specific inFlightPacket in the store from its index
func (k Keeper) SetOutgoingInFlightPacket(ctx context.Context, inFlightPacket types.OutgoingInFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OutgoingInFlightPacketKeyPrefix))
	b := k.cdc.MustMarshal(&inFlightPacket)
	store.Set(types.OutgoingInFlightPacketKey(inFlightPacket.Index), b)
}

// GetOutgoingInFlightPacket returns a inFlightPacket from its index
func (k Keeper) GetOutgoingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) (val types.OutgoingInFlightPacket, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OutgoingInFlightPacketKeyPrefix))

	b := store.Get(types.OutgoingInFlightPacketKey(
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

// RemoveOutgoingInFlightPacket removes a inFlightPacket from the store
func (k Keeper) RemoveOutgoingInFlightPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OutgoingInFlightPacketKeyPrefix))
	store.Delete(types.OutgoingInFlightPacketKey(
		types.NewPacketIndex(
			srcPortId,
			srcChannelId,
			sequence,
		),
	))
}

// GetAllOutgoingInFlightPacket returns all inFlightPacket
func (k Keeper) GetAllOutgoingInFlightPacket(ctx context.Context) (list []types.OutgoingInFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.OutgoingInFlightPacketKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OutgoingInFlightPacket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
