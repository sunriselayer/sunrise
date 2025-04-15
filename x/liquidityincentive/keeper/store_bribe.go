package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SetBribe set a specific bribe in the store from its index
func (k Keeper) SetBribe(ctx context.Context, bribe types.Bribe) error {
	err := k.Bribes.Set(ctx, types.BribeKey(bribe.EpochId, bribe.PoolId), bribe)
	if err != nil {
		return err
	}
	return nil
}

// GetBribe returns a bribe from its index
func (k Keeper) GetBribe(ctx context.Context, epochId uint64, poolId uint64) (val types.Bribe, found bool, err error) {
	key := types.BribeKey(epochId, poolId)
	has, err := k.Bribes.Has(ctx, key)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Bribes.Get(ctx, key)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveBribe removes a bribe from the store
func (k Keeper) RemoveBribe(ctx context.Context, epochId uint64, poolId uint64) error {
	err := k.Bribes.Remove(ctx, types.BribeKey(epochId, poolId))
	if err != nil {
		return err
	}
	return nil
}

// GetAllBribes returns all bribes
func (k Keeper) GetAllBribes(ctx context.Context) (list []types.Bribe, err error) {
	err = k.Bribes.Walk(
		ctx,
		nil,
		func(key collections.Pair[uint64, uint64], value types.Bribe) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAllBribeByEpochId returns all bribes by epoch id
func (k Keeper) GetAllBribeByEpochId(ctx context.Context, epochId uint64) (list []types.Bribe, err error) {
	err = k.Bribes.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](epochId),
		func(key collections.Pair[uint64, uint64], value types.Bribe) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
