package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetPublishedData(ctx context.Context, dataHash []byte) (data types.PublishedData) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.PublishedDataKey(dataHash))
	if bz == nil {
		return data
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data
}

// SetParams set the params
func (k Keeper) SetPublishedData(ctx context.Context, data types.PublishedData) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, err := k.cdc.Marshal(&data)
	if err != nil {
		return err
	}
	store.Set(types.PublishedDataKey(data.RecoveredDataHash), bz)
	return nil
}

func (k Keeper) DeletePublishedData(ctx sdk.Context, dataHash []byte) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.PublishedDataKey(dataHash))
}

func (k Keeper) GetAllPublishedData(ctx sdk.Context) []types.PublishedData {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.PublishedDataKeyPrefix)
	defer iter.Close()

	data := []types.PublishedData{}
	for ; iter.Valid(); iter.Next() {
		da := types.PublishedData{}
		k.cdc.MustUnmarshal(iter.Value(), &da)
		data = append(data, da)
	}
	return data
}
