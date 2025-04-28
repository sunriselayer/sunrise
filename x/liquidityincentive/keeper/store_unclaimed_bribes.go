package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// GetUnclaimedBribe returns the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) GetUnclaimedBribe(ctx context.Context, voter sdk.AccAddress, bribeId uint64) (types.UnclaimedBribe, error) {
	key := collections.Join(voter, bribeId)
	return k.UnclaimedBribes.Get(ctx, key)
}

// SetUnclaimedBribe sets the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) SetUnclaimedBribe(ctx context.Context, unclaimedBribe types.UnclaimedBribe) error {
	address, err := k.addressCodec.StringToBytes(unclaimedBribe.Address)
	if err != nil {
		return err
	}
	acc := sdk.AccAddress(address)
	key := collections.Join(acc, unclaimedBribe.BribeId)
	return k.UnclaimedBribes.Set(ctx, key, unclaimedBribe)
}

// RemoveUnclaimedBribe removes the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) RemoveUnclaimedBribe(ctx context.Context, unclaimedBribe types.UnclaimedBribe) error {
	address, err := k.addressCodec.StringToBytes(unclaimedBribe.Address)
	if err != nil {
		return err
	}
	acc := sdk.AccAddress(address)
	key := collections.Join(acc, unclaimedBribe.BribeId)
	return k.UnclaimedBribes.Remove(ctx, key)
}

// GetAllUnclaimedBribes returns all unclaimed bribes
func (k Keeper) GetAllUnclaimedBribes(ctx context.Context) (list []types.UnclaimedBribe, err error) {
	err = k.UnclaimedBribes.Walk(
		ctx,
		nil,
		func(key collections.Pair[sdk.AccAddress, uint64], value types.UnclaimedBribe) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}
