package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// AppendBribe appends a bribe in the store with a new id and update the count
func (k Keeper) AppendBribe(ctx context.Context, bribe types.Bribe) (id uint64, err error) {
	id, err = k.BribeId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	bribe.Id = id
	err = k.SetBribe(ctx, bribe)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetBribe set a specific bribe in the store from its index
func (k Keeper) SetBribe(ctx context.Context, bribe types.Bribe) error {
	err := k.Bribes.Set(ctx, bribe.Id, bribe)
	if err != nil {
		return err
	}
	return nil
}

// GetBribe returns a bribe from its index
func (k Keeper) GetBribe(ctx context.Context, id uint64) (val types.Bribe, found bool, err error) {
	has, err := k.Bribes.Has(ctx, id)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Bribes.Get(ctx, id)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveBribe removes a bribe from the store
func (k Keeper) RemoveBribe(ctx context.Context, id uint64) error {
	err := k.Bribes.Remove(ctx, id)
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
		func(key uint64, value types.Bribe) (bool, error) {
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
	err = k.Bribes.Indexes.EpochId.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](epochId),
		func(_ uint64, bribeId uint64) (bool, error) {
			bribe, _, err := k.GetBribe(ctx, bribeId)
			if err != nil {
				return false, err
			}
			list = append(list, bribe)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAllBribeByPoolId returns all bribes by pool id
func (k Keeper) GetAllBribeByPoolId(ctx context.Context, poolId uint64) (list []types.Bribe, err error) {
	err = k.Bribes.Indexes.PoolId.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](poolId),
		func(_ uint64, bribeId uint64) (bool, error) {
			bribe, _, err := k.GetBribe(ctx, bribeId)
			if err != nil {
				return false, err
			}
			list = append(list, bribe)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetBribeByEpochAndPool retrieves a bribe by epoch ID and pool ID.
// Note: This iterates through all bribes, which might be inefficient for a large number of bribes.
func (k Keeper) GetBribeByEpochAndPool(ctx context.Context, epochId, poolId uint64) (types.Bribe, bool, error) {
	var foundBribe types.Bribe
	found := false
	err := k.Bribes.Walk(ctx, nil, func(key uint64, bribe types.Bribe) (stop bool, err error) {
		if bribe.EpochId == epochId && bribe.PoolId == poolId {
			foundBribe = bribe
			found = true
			return true, nil // Stop iteration once found
		}
		return false, nil
	})
	if err != nil {
		return types.Bribe{}, false, err
	}
	return foundBribe, found, nil
}
