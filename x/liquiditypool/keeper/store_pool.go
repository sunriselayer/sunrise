package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// GetPoolCount get the total number of pool
func (k Keeper) GetPoolCount(ctx context.Context) (count uint64, err error) {
	count, err = k.PoolId.Peek(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SetPoolCount set the total number of pool
func (k Keeper) SetPoolCount(ctx context.Context, count uint64) error {
	err := k.PoolId.Set(ctx, count)
	if err != nil {
		return err
	}
	return nil
}

// AppendPool appends a pool in the store with a new id and update the count
func (k Keeper) AppendPool(ctx context.Context, pool types.Pool) (id uint64, err error) {
	// Create the pool
	id, err = k.PoolId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	pool.Id = id
	err = k.SetPool(ctx, pool)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetPool set a specific pool in the store
func (k Keeper) SetPool(ctx context.Context, pool types.Pool) error {
	err := k.Pools.Set(ctx, pool.Id, pool)
	if err != nil {
		return err
	}
	return nil
}

// GetPool returns a pool from its id
func (k Keeper) GetPool(ctx context.Context, id uint64) (val types.Pool, found bool, err error) {
	has, err := k.Pools.Has(ctx, id)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Pools.Get(ctx, id)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx context.Context, id uint64) error {
	err := k.Pools.Remove(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllPools returns all pool
func (k Keeper) GetAllPools(ctx context.Context) (list []types.Pool, err error) {
	err = k.Pools.Walk(
		ctx,
		nil,
		func(key uint64, value types.Pool) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return
}
