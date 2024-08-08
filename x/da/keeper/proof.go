package keeper

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetProof(ctx context.Context, metadataUri string, sender string) (data types.Proof) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.ProofKey(metadataUri, sender))
	if bz == nil {
		return data
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data
}

// SetParams set the params
func (k Keeper) SetProof(ctx context.Context, data types.Proof) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, err := k.cdc.Marshal(&data)
	if err != nil {
		return err
	}
	store.Set(types.ProofKey(data.MetadataUri, data.Sender), bz)
	return nil
}

func (k Keeper) DeleteProof(ctx sdk.Context, metadataUri string, sender string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.ProofKey(metadataUri, sender))
}

func (k Keeper) GetAllProofs(ctx sdk.Context) []types.Proof {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.ProofKeyPrefix)
	defer iter.Close()

	data := []types.Proof{}
	for ; iter.Valid(); iter.Next() {
		da := types.Proof{}
		k.cdc.MustUnmarshal(iter.Value(), &da)
		data = append(data, da)
	}
	return data
}
