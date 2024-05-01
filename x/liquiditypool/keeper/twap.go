package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// SetTwap set a specific twap in the store from its index
func (k Keeper) SetTwap(ctx context.Context, twap types.Twap) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store :=  prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))
	b := k.cdc.MustMarshal(&twap)
	store.Set(types.TwapKey(
        twap.Index,
    ), b)
}

// GetTwap returns a twap from its index
func (k Keeper) GetTwap(
    ctx context.Context,
    index string,
    
) (val types.Twap, found bool) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))

	b := store.Get(types.TwapKey(
        index,
    ))
    if b == nil {
        return val, false
    }

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTwap removes a twap from the store
func (k Keeper) RemoveTwap(
    ctx context.Context,
    index string,
    
) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))
	store.Delete(types.TwapKey(
	    index,
    ))
}

// GetAllTwap returns all twap
func (k Keeper) GetAllTwap(ctx context.Context) (list []types.Twap) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TwapKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Twap
		k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
	}

    return
}
