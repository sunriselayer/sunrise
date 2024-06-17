package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// GetEpochCount get the total number of epoch
func (k Keeper) GetEpochCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.EpochCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetEpochCount set the total number of epoch
func (k Keeper) SetEpochCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.EpochCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendEpoch appends a epoch in the store with a new id and update the count
func (k Keeper) AppendEpoch(ctx context.Context, epoch types.Epoch) uint64 {
	// Create the epoch
	count := k.GetEpochCount(ctx)

	// Set the ID of the appended value
	epoch.Id = count

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EpochKey))
	appendedValue := k.cdc.MustMarshal(&epoch)
	store.Set(GetEpochIDBytes(epoch.Id), appendedValue)

	// Update epoch count
	k.SetEpochCount(ctx, count+1)

	return count
}

// SetEpoch set a specific epoch in the store
func (k Keeper) SetEpoch(ctx context.Context, epoch types.Epoch) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EpochKey))
	b := k.cdc.MustMarshal(&epoch)
	store.Set(GetEpochIDBytes(epoch.Id), b)
}

// GetEpoch returns a epoch from its id
func (k Keeper) GetEpoch(ctx context.Context, id uint64) (val types.Epoch, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EpochKey))
	b := store.Get(GetEpochIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveEpoch removes a epoch from the store
func (k Keeper) RemoveEpoch(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EpochKey))
	store.Delete(GetEpochIDBytes(id))
}

func (k Keeper) GetLastEpoch(ctx context.Context) (epoch types.Epoch, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EpochKey))
	iterator := storetypes.KVStoreReversePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshal(iterator.Value(), &epoch)
		return epoch, true
	}

	return types.Epoch{}, false
}

// GetAllEpoch returns all epoch
func (k Keeper) GetAllEpoch(ctx context.Context) (list []types.Epoch) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.EpochKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Epoch
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetEpochIDBytes returns the byte representation of the ID
func GetEpochIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.EpochKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}
