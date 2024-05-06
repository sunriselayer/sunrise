package keeper

import (
	"context"
	"time"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// SetTradeFootprint set a specific TradeFootprint in the store from its index
func (k Keeper) SetTradeFootprint(ctx context.Context, baseDenom string, quoteDenom string, tradeFootprint types.TradeFootprint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeFootprintKeyPrefix))
	b := k.cdc.MustMarshal(&tradeFootprint)
	store.Set(types.TradeFootprintKey(
		baseDenom,
		quoteDenom,
		tradeFootprint.Timestamp,
	), b)
}

// GetTradeFootprint returns a TradeFootprint from its index
func (k Keeper) GetTradeFootprint(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
	timestamp time.Time,
) (val types.TradeFootprint, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeFootprintKeyPrefix))

	b := store.Get(types.TradeFootprintKey(
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

// RemoveTradeFootprint removes a TradeFootprint from the store
func (k Keeper) RemoveTradeFootprint(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
	timestamp time.Time,
) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeFootprintKeyPrefix))
	store.Delete(types.TradeFootprintKey(
		baseDenom,
		quoteDenom,
		timestamp,
	))
}

// GetAllTradeFootprint returns all TradeFootprint
func (k Keeper) GetAllTradeFootprint(
	ctx context.Context,
	baseDenom string,
	quoteDenom string,
) (list []types.TradeFootprint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TradeFootprintKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, types.PriceFootprintIterationPrefix(baseDenom, quoteDenom))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TradeFootprint
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
