package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/sunriselayer/sunrise/x/blobgrant/types"
)

// SetRegistration set a specific registration in the store from its index
func (k Keeper) SetRegistration(ctx context.Context, registration types.Registration) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RegistrationKeyPrefix))
	b := k.cdc.MustMarshal(&registration)
	store.Set(types.RegistrationKey(
		registration.LiquidityProvider,
	), b)
}

// GetRegistration returns a registration from its index
func (k Keeper) GetRegistration(
	ctx context.Context,
	address string,

) (val types.Registration, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RegistrationKeyPrefix))

	b := store.Get(types.RegistrationKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveRegistration removes a registration from the store
func (k Keeper) RemoveRegistration(
	ctx context.Context,
	address string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RegistrationKeyPrefix))
	store.Delete(types.RegistrationKey(
		address,
	))
}

// GetAllRegistration returns all registration
func (k Keeper) GetAllRegistration(ctx context.Context) (list []types.Registration) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.RegistrationKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Registration
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
