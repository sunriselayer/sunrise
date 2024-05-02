package keeper

import (
	"context"
	"time"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// SetPriceFootprint set a specific PriceFootprint in the store from its index
func (k Keeper) SetPriceFootprint(ctx context.Context, baseDenom string, quoteDenom string, priceFootprint types.PriceFootprint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PriceFootprintKeyPrefix))
	b := k.cdc.MustMarshal(&priceFootprint)
	store.Set(types.PriceFootprintKey(
		baseDenom,
		quoteDenom,
		priceFootprint.Timestamp,
	), b)
}

// GetPriceFootprint returns a PriceFootprint from its index
func (k Keeper) GetPriceFootprint(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
	timestamp time.Time,
) (val types.PriceFootprint, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PriceFootprintKeyPrefix))

	b := store.Get(types.PriceFootprintKey(
		baseDenom,
		quoteDenom,
		timestamp,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePriceFootprint removes a PriceFootprint from the store
func (k Keeper) RemovePriceFootprint(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
	timestamp time.Time,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PriceFootprintKeyPrefix))
	store.Delete(types.PriceFootprintKey(
		baseDenom,
		quoteDenom,
		timestamp,
	))
}

// GetAllPriceFootprint returns all PriceFootprint
func (k Keeper) GetAllPriceFootprint(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
) (list []types.PriceFootprint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.PriceFootprintKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, types.PriceFootprintIterationPrefix(baseDenom, quoteDenom))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PriceFootprint
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
