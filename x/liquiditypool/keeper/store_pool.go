package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// GetPoolCount get the total number of pool
func (k Keeper) GetPoolCount(ctx context.Context) (uint64, error) {
	return k.PoolId.Peek(ctx)
}

// SetPoolCount set the total number of pool
func (k Keeper) SetPoolCount(ctx context.Context, count uint64) error {
	return k.PoolId.Set(ctx, count)
}

// AppendPool appends a pool in the store with a new id and update the count
func (k Keeper) AppendPool(ctx context.Context, pool types.Pool) (uint64, error) {
	// Create the pool
	id, err := k.PoolId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	pool.Id = id
	if err := k.SetPool(ctx, pool); err != nil {
		return 0, err
	}

	return id, nil
}

// SetPool set a specific pool in the store
func (k Keeper) SetPool(ctx context.Context, pool types.Pool) error {
	return k.Pools.Set(ctx, pool.Id, pool)
}

// GetPool returns the pool for the given id
func (k Keeper) GetPool(ctx context.Context, id uint64) (pool types.Pool, found bool, err error) {
	has, err := k.Pools.Has(ctx, id)
	if err != nil {
		return pool, false, err
	}
	if !has {
		return pool, false, nil
	}
	val, err := k.Pools.Get(ctx, id)
	if err != nil {
		return pool, false, err
	}
	return val, true, nil
}

// DeletePool removes the pool
func (k Keeper) DeletePool(ctx context.Context, id uint64) error {
	return k.Pools.Remove(ctx, id)
}

// GetAllPools returns all pools
func (k Keeper) GetAllPools(ctx context.Context) (list []types.Pool, err error) {
	err = k.Pools.Walk(ctx, nil, func(key uint64, value types.Pool) (bool, error) {
		list = append(list, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
