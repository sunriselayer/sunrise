package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// GetUnclaimedBribe returns the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) GetUnclaimedBribe(ctx context.Context, voter sdk.AccAddress, epochId uint64, poolId uint64) (types.UnclaimedBribe, error) {
	key := collections.Join3(voter, epochId, poolId)
	return k.UnclaimedBribes.Get(ctx, key)
}

// SetUnclaimedBribe sets the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) SetUnclaimedBribe(ctx context.Context, voter sdk.AccAddress, epochId uint64, poolId uint64, unclaimedBribe types.UnclaimedBribe) error {
	key := collections.Join3(voter, epochId, poolId)
	return k.UnclaimedBribes.Set(ctx, key, unclaimedBribe)
}

// RemoveUnclaimedBribe removes the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) RemoveUnclaimedBribe(ctx context.Context, voter sdk.AccAddress, epochId uint64, poolId uint64) error {
	key := collections.Join3(voter, epochId, poolId)
	return k.UnclaimedBribes.Remove(ctx, key)
}

// GetAllUnclaimedBribes returns all unclaimed bribes
func (k Keeper) GetAllUnclaimedBribes(ctx context.Context) ([]types.UnclaimedBribe, error) {
	var unclaimedBribes []types.UnclaimedBribe
	err := k.UnclaimedBribes.Walk(ctx, nil, func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.UnclaimedBribe) (bool, error) {
		unclaimedBribes = append(unclaimedBribes, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return unclaimedBribes, nil
}

// GetUnclaimedBribesByEpoch returns all unclaimed bribes for a given epoch
func (k Keeper) GetUnclaimedBribesByEpoch(ctx context.Context, epochId uint64) ([]types.UnclaimedBribe, error) {
	var unclaimedBribes []types.UnclaimedBribe
	err := k.UnclaimedBribes.Walk(ctx, collections.NewPrefixedTripleRange[sdk.AccAddress, uint64, uint64](sdk.AccAddress{}),
		func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.UnclaimedBribe) (bool, error) {
			if key.K2() == epochId {
				unclaimedBribes = append(unclaimedBribes, value)
			}
			return false, nil
		})
	if err != nil {
		return nil, err
	}
	return unclaimedBribes, nil
}

// GetUnclaimedBribesByVoter returns all unclaimed bribes for a given voter
func (k Keeper) GetUnclaimedBribesByVoter(ctx context.Context, voter sdk.AccAddress) ([]types.UnclaimedBribe, error) {
	var unclaimedBribes []types.UnclaimedBribe
	err := k.UnclaimedBribes.Walk(ctx, collections.NewPrefixedTripleRange[sdk.AccAddress, uint64, uint64](voter),
		func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.UnclaimedBribe) (bool, error) {
			unclaimedBribes = append(unclaimedBribes, value)
			return false, nil
		})
	if err != nil {
		return nil, err
	}
	return unclaimedBribes, nil
}
