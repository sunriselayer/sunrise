package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetInFlightPacket set a specific inFlightPacket in the store from its index
func (k Keeper) SetInFlightPacket(ctx context.Context, inFlightPacket types.InFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.InFlightPacketKeyPrefix))
	b := k.cdc.MustMarshal(&inFlightPacket)
	store.Set(types.InFlightPacketKey(
		inFlightPacket.Index,
	), b)
}

// GetInFlightPacket returns a inFlightPacket from its index
func (k Keeper) GetInFlightPacket(
	ctx context.Context,
	index string,

) (val types.InFlightPacket, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.InFlightPacketKeyPrefix))

	b := store.Get(types.InFlightPacketKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveInFlightPacket removes a inFlightPacket from the store
func (k Keeper) RemoveInFlightPacket(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.InFlightPacketKeyPrefix))
	store.Delete(types.InFlightPacketKey(
		index,
	))
}

// GetAllInFlightPacket returns all inFlightPacket
func (k Keeper) GetAllInFlightPacket(ctx context.Context) (list []types.InFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.InFlightPacketKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InFlightPacket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
