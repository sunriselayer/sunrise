package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetPublishedData(ctx context.Context, metadataUri string) (data types.PublishedData) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.PublishedDataKey(metadataUri))
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
	store.Set(types.PublishedDataKey(data.MetadataUri), bz)

	if data.Status == "verified" {
		store.Delete(types.UnverifiedDataByTimeKey(uint64(data.Timestamp.Unix()), data.MetadataUri))
	} else {
		store.Set(types.UnverifiedDataByTimeKey(uint64(data.Timestamp.Unix()), data.MetadataUri), []byte(data.MetadataUri))
	}
	return nil
}

func (k Keeper) DeletePublishedData(ctx sdk.Context, data types.PublishedData) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.PublishedDataKey(data.MetadataUri))
	store.Delete(types.UnverifiedDataByTimeKey(uint64(data.Timestamp.Unix()), data.MetadataUri))
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

func (k Keeper) GetUnverifiedDataBeforeTime(ctx sdk.Context, timestamp uint64) ([]types.PublishedData, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.UnverifiedDataByTimeKeyPrefix)

	iterator := prefixStore.Iterator(nil, sdk.Uint64ToBigEndian(timestamp))
	defer iterator.Close()

	dataArray := []types.PublishedData{}
	for ; iterator.Valid(); iterator.Next() {
		metadataUri := string(iterator.Value())
		data := k.GetPublishedData(ctx, metadataUri)
		dataArray = append(dataArray, data)
	}
	return dataArray, nil
}
