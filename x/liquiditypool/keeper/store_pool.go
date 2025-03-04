package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// GetPoolCount get the total number of pool
func (k Keeper) GetPoolCount(ctx context.Context) uint64 {
	val, err := k.PoolId.Peek(ctx)
	if err != nil {
		panic(err)
	}

	return val
}

// SetPoolCount set the total number of pool
func (k Keeper) SetPoolCount(ctx context.Context, count uint64) {
	err := k.PoolId.Set(ctx, count)
	if err != nil {
		panic(err)
	}
}

// AppendPool appends a pool in the store with a new id and update the count
func (k Keeper) AppendPool(ctx context.Context, pool types.Pool) uint64 {
	// Create the pool
	id, err := k.PoolId.Next(ctx)
	if err != nil {
		panic(err)
	}

	// Set the ID of the appended value
	pool.Id = id
	k.SetPool(ctx, pool)

	return id
}

// SetPool set a specific pool in the store
func (k Keeper) SetPool(ctx context.Context, pool types.Pool) {
	err := k.Pools.Set(ctx, pool.Id, pool)
	if err != nil {
		panic(err)
	}
}

// GetPool returns a pool from its id
func (k Keeper) GetPool(ctx context.Context, id uint64) (val types.Pool, found bool) {
	has, err := k.Pools.Has(ctx, id)
	if err != nil {
		panic(err)
	}

	if !has {
		return val, false
	}

	val, err = k.Pools.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx context.Context, id uint64) {
	err := k.Pools.Remove(ctx, id)
	if err != nil {
		panic(err)
	}
}

// GetAllPools returns all pool
func (k Keeper) GetAllPools(ctx context.Context) (list []types.Pool) {
	err := k.Pools.Walk(
		ctx,
		nil,
		func(key uint64, value types.Pool) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}
