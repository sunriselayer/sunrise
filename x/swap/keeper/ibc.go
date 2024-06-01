package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	types "github.com/sunriselayer/sunrise/x/swap/types"
)

func (k Keeper) SetInFlightPackets(ctx context.Context, inFlightPackets map[string]types.InFlightPacket) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.InFlightPacketKey))

	for key, value := range inFlightPackets {
		key := key
		value := value
		bz := k.cdc.MustMarshal(&value)
		store.Set([]byte(key), bz)
	}
}

func (k Keeper) GetInFlightPackets(ctx context.Context) map[string]types.InFlightPacket {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.InFlightPacketKey))

	inFlightPackets := make(map[string]types.InFlightPacket)

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InFlightPacket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		inFlightPackets[string(iterator.Key())] = val
	}

	return inFlightPackets
}
