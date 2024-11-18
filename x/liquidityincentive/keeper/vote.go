package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SetVote set a specific vote in the store from its index
func (k Keeper) SetVote(ctx context.Context, vote types.Vote) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.VoteKeyPrefix))
	b := k.cdc.MustMarshal(&vote)
	store.Set([]byte(vote.Sender), b)

}

// GetVote returns a vote from its index
func (k Keeper) GetVote(ctx context.Context, sender string) (val types.Vote, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.VoteKeyPrefix))

	b := store.Get([]byte(sender))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveVote removes a vote from the store
func (k Keeper) RemoveVote(ctx context.Context, sender string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.VoteKeyPrefix))
	store.Delete([]byte(sender))
}

// GetAllVote returns all vote
func (k Keeper) GetAllVotes(ctx context.Context) (list []types.Vote) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.VoteKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
