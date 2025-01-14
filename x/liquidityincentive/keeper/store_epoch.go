package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// GetEpochCount get the total number of epoch
func (k Keeper) GetEpochCount(ctx context.Context) uint64 {
	val, err := k.EpochId.Peek(ctx)
	if err != nil {
		panic(err)
	}

	return val
}

// SetEpochCount set the total number of epoch
func (k Keeper) SetEpochCount(ctx context.Context, count uint64) {
	err := k.EpochId.Set(ctx, count)
	if err != nil {
		panic(err)
	}
}

// AppendEpoch appends a epoch in the store with a new id and update the count
func (k Keeper) AppendEpoch(ctx context.Context, epoch types.Epoch) uint64 {
	// Create the epoch
	id, err := k.EpochId.Next(ctx)
	if err != nil {
		panic(err)
	}

	// Set the ID of the appended value
	epoch.Id = id

	k.SetEpoch(ctx, epoch)

	return id
}

// SetEpoch set a specific epoch in the store
func (k Keeper) SetEpoch(ctx context.Context, epoch types.Epoch) {
	err := k.Epochs.Set(ctx, epoch.Id, epoch)
	if err != nil {
		panic(err)
	}
}

// GetEpoch returns a epoch from its id
func (k Keeper) GetEpoch(ctx context.Context, id uint64) (val types.Epoch, found bool) {
	has, err := k.Epochs.Has(ctx, id)
	if err != nil {
		panic(err)
	}

	if !has {
		return val, false
	}

	val, err = k.Epochs.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	return val, true
}

// RemoveEpoch removes a epoch from the store
func (k Keeper) RemoveEpoch(ctx context.Context, id uint64) {
	err := k.Epochs.Remove(ctx, id)
	if err != nil {
		panic(err)
	}
}

// GetAllEpoch returns all epoch
func (k Keeper) GetAllEpoch(ctx context.Context) (list []types.Epoch) {
	err := k.Epochs.Walk(
		ctx,
		nil,
		func(key uint64, value types.Epoch) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		panic(err)
	}

	return
}

func (k Keeper) GetLastEpoch(ctx context.Context) (epoch types.Epoch, found bool) {
	has, err := k.Epochs.Has(ctx, 0)
	if err != nil {
		panic(err)
	}

	if !has {
		return epoch, false
	}

	err = k.Epochs.Walk(
		ctx,
		new(collections.Range[uint64]).Descending(),
		func(key uint64, value types.Epoch) (bool, error) {
			epoch = value
			return true, nil
		},
	)

	if err != nil {
		panic(err)
	}

	return epoch, true
}
