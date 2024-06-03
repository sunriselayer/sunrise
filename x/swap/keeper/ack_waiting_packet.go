package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// SetAckWaitingPacket set a specific ackWaitingPacket in the store from its index
func (k Keeper) SetAckWaitingPacket(ctx context.Context, ackWaitingPacket types.AckWaitingPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AckWaitingPacketKeyPrefix))
	b := k.cdc.MustMarshal(&ackWaitingPacket)
	store.Set(types.AckWaitingPacketKey(
		ackWaitingPacket.Index,
	), b)
}

// GetAckWaitingPacket returns a ackWaitingPacket from its index
func (k Keeper) GetAckWaitingPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) (val types.AckWaitingPacket, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AckWaitingPacketKeyPrefix))

	b := store.Get(types.AckWaitingPacketKey(
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

// RemoveAckWaitingPacket removes a ackWaitingPacket from the store
func (k Keeper) RemoveAckWaitingPacket(
	ctx context.Context,
	srcPortId string,
	srcChannelId string,
	sequence uint64,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AckWaitingPacketKeyPrefix))
	store.Delete(types.AckWaitingPacketKey(
		types.NewPacketIndex(
			srcPortId,
			srcChannelId,
			sequence,
		),
	))
}

// GetAllAckWaitingPacket returns all ackWaitingPacket
func (k Keeper) GetAllAckWaitingPacket(ctx context.Context) (list []types.AckWaitingPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.AckWaitingPacketKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AckWaitingPacket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
