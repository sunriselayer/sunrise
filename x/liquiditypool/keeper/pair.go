package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// SetPair set a specific pair in the store from its index
func (k Keeper) SetPair(ctx context.Context, pair types.Pair) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store :=  prefix.NewStore(storeAdapter, types.KeyPrefix(types.PairKeyPrefix))
	b := k.cdc.MustMarshal(&pair)
	store.Set(types.PairKey(
        pair.Index,
    ), b)
}

// GetPair returns a pair from its index
func (k Keeper) GetPair(
    ctx context.Context,
    index string,
    
) (val types.Pair, found bool) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PairKeyPrefix))

	b := store.Get(types.PairKey(
        index,
    ))
    if b == nil {
        return val, false
    }

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePair removes a pair from the store
func (k Keeper) RemovePair(
    ctx context.Context,
    index string,
    
) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PairKeyPrefix))
	store.Delete(types.PairKey(
	    index,
    ))
}

// GetAllPair returns all pair
func (k Keeper) GetAllPair(ctx context.Context) (list []types.Pair) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PairKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pair
		k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
	}

    return
}
