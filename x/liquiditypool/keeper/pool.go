package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// GetPoolCount get the total number of pool
func (k Keeper) GetPoolCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.PoolCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetPoolCount set the total number of pool
func (k Keeper) SetPoolCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.PoolCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPool appends a pool in the store with a new id and update the count
func (k Keeper) AppendPool(
	ctx context.Context,
	pool types.Pool,
) uint64 {
	// Create the pool
	count := k.GetPoolCount(ctx)

	// Set the ID of the appended value
	pool.Id = count

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PoolKey))
	appendedValue := k.cdc.MustMarshal(&pool)
	store.Set(GetPoolIDBytes(pool.Id), appendedValue)

	// Update pool count
	k.SetPoolCount(ctx, count+1)

	return count
}

// SetPool set a specific pool in the store
func (k Keeper) SetPool(ctx context.Context, pool types.Pool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PoolKey))
	b := k.cdc.MustMarshal(&pool)
	store.Set(GetPoolIDBytes(pool.Id), b)
}

// GetPool returns a pool from its id
func (k Keeper) GetPool(ctx context.Context, id uint64) (val types.Pool, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PoolKey))
	b := store.Get(GetPoolIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PoolKey))
	store.Delete(GetPoolIDBytes(id))
}

// GetAllPool returns all pool
func (k Keeper) GetAllPool(ctx context.Context) (list []types.Pool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PoolKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPoolIDBytes returns the byte representation of the ID
func GetPoolIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.PoolKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}

func (k Keeper) GetPoolBalance(ctx context.Context, pool types.Pool) (x math.Int, y math.Int) {
	addr := k.accountKeeper.GetModuleAddress(types.PoolModuleName(pool.Id))
	balances := k.bankKeeper.SpendableCoins(ctx, addr)

	for _, balance := range balances {
		if balance.Denom == pool.BaseDenom {
			x = balance.Amount
		} else if balance.Denom == pool.QuoteDenom {
			y = balance.Amount
		}
	}

	return
}

func (k Keeper) GetLpTokenSupply(ctx context.Context, poolId uint64) sdk.Coin {
	return k.bankKeeper.GetSupply(ctx, types.LpTokenDenom(poolId))
}
