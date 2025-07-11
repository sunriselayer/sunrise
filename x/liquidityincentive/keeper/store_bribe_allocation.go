package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// GetBribeAllocation returns the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) GetBribeAllocation(ctx context.Context, voter sdk.AccAddress, epochId uint64, poolId uint64) (types.BribeAllocation, error) {
	key := collections.Join3(voter, epochId, poolId)
	return k.BribeAllocations.Get(ctx, key)
}

// SetBribeAllocation sets the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) SetBribeAllocation(ctx context.Context, unclaimedBribe types.BribeAllocation) error {
	address, err := k.addressCodec.StringToBytes(unclaimedBribe.Address)
	if err != nil {
		return err
	}
	acc := sdk.AccAddress(address)
	key := collections.Join3(acc, unclaimedBribe.EpochId, unclaimedBribe.PoolId)
	return k.BribeAllocations.Set(ctx, key, unclaimedBribe)
}

// RemoveBribeAllocation removes the unclaimed bribe for a given voter, epoch and pool
func (k Keeper) RemoveBribeAllocation(ctx context.Context, key collections.Triple[sdk.AccAddress, uint64, uint64]) error {
	return k.BribeAllocations.Remove(ctx, key)
}

// GetAllBribeAllocations returns all unclaimed bribes
func (k Keeper) GetAllBribeAllocations(ctx context.Context) (list []types.BribeAllocation, err error) {
	err = k.BribeAllocations.Walk(
		ctx,
		nil,
		func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.BribeAllocation) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetBribeAllocationsByAddress returns all bribe allocations by address
func (k Keeper) GetBribeAllocationsByAddress(ctx context.Context, address sdk.AccAddress) (list []types.BribeAllocation, err error) {
	prefix := collections.NewPrefixedTripleRange[sdk.AccAddress, uint64, uint64](address)
	err = k.BribeAllocations.Walk(ctx, prefix, func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.BribeAllocation) (stop bool, err error) {
		list = append(list, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetBribeAllocationsByAddressAndEpochId returns all bribe allocations by address and epoch id
func (k Keeper) GetBribeAllocationsByAddressAndEpochId(ctx context.Context, address sdk.AccAddress, epochId uint64) (list []types.BribeAllocation, err error) {
	prefix := collections.NewSuperPrefixedTripleRange[sdk.AccAddress, uint64, uint64](address, epochId)
	err = k.BribeAllocations.Walk(ctx, prefix, func(key collections.Triple[sdk.AccAddress, uint64, uint64], value types.BribeAllocation) (stop bool, err error) {
		list = append(list, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
