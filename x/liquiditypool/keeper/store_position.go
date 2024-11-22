package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetPositionCount get the total number of position
func (k Keeper) GetPositionCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.PositionCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetPositionCount set the total number of position
func (k Keeper) SetPositionCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.PositionCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPosition appends a position in the store with a new id and update the count
func (k Keeper) AppendPosition(ctx context.Context, position types.Position) uint64 {
	// Create the position
	count := k.GetPositionCount(ctx)

	// Set the ID of the appended value
	position.Id = count
	k.SetPosition(ctx, position)

	// Update position count
	k.SetPositionCount(ctx, count+1)

	return count
}

// SetPosition set a specific position in the store
func (k Keeper) SetPosition(ctx context.Context, position types.Position) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PositionKey))
	b := k.cdc.MustMarshal(&position)
	store.Set(GetPositionIDBytes(position.Id), b)

	positionKey := sdk.Uint64ToBigEndian(position.Id)
	store = prefix.NewStore(storeAdapter, types.PositionByPoolPrefix(position.PoolId))
	store.Set(positionKey, positionKey)

	store = prefix.NewStore(storeAdapter, types.PositionByAddressPrefix(position.Address))
	store.Set(positionKey, positionKey)
}

// GetPosition returns a position from its id
func (k Keeper) GetPosition(ctx context.Context, id uint64) (val types.Position, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PositionKey))
	b := store.Get(GetPositionIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePosition removes a position from the store
func (k Keeper) RemovePosition(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PositionKey))
	store.Delete(GetPositionIDBytes(id))
}

// GetAllPositions returns all position
func (k Keeper) GetAllPositions(ctx context.Context) (list []types.Position) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PositionKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPositionIDBytes returns the byte representation of the ID
func GetPositionIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.PositionKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}

func (k Keeper) PoolHasPosition(ctx context.Context, poolId uint64) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PositionByPoolPrefix(poolId))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		return true
	}
	return false
}

func (k Keeper) GetPositionsByPool(ctx context.Context, poolId uint64) []types.Position {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PositionByPoolPrefix(poolId))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	positions := []types.Position{}
	for ; iterator.Valid(); iterator.Next() {
		id := sdk.BigEndianToUint64(iterator.Value())
		position, found := k.GetPosition(ctx, id)
		if found {
			positions = append(positions, position)
		}
	}
	return positions
}

func (k Keeper) GetPositionsByAddress(ctx context.Context, addr string) []types.Position {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PositionByAddressPrefix(addr))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	positions := []types.Position{}
	for ; iterator.Valid(); iterator.Next() {
		id := sdk.BigEndianToUint64(iterator.Value())
		position, found := k.GetPosition(ctx, id)
		if found {
			positions = append(positions, position)
		}
	}
	return positions
}
