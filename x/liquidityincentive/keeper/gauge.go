package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SetGauge set a specific gauge in the store from its index
func (k Keeper) SetGauge(ctx context.Context, gauge types.Gauge) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GaugeKeyPrefix))
	b := k.cdc.MustMarshal(&gauge)
	store.Set(types.GaugeKey(
		gauge.Index,
	), b)
}

// GetGauge returns a gauge from its index
func (k Keeper) GetGauge(
	ctx context.Context,
	index string,

) (val types.Gauge, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GaugeKeyPrefix))

	b := store.Get(types.GaugeKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGauge removes a gauge from the store
func (k Keeper) RemoveGauge(
	ctx context.Context,
	index string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GaugeKeyPrefix))
	store.Delete(types.GaugeKey(
		index,
	))
}

// GetAllGauge returns all gauge
func (k Keeper) GetAllGauge(ctx context.Context) (list []types.Gauge) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.GaugeKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Gauge
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
