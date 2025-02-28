package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// GetEpochCount get the total number of epoch
func (k Keeper) GetEpochCount(ctx context.Context) (uint64, error) {
	val, err := k.EpochId.Peek(ctx)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// SetEpochCount set the total number of epoch
func (k Keeper) SetEpochCount(ctx context.Context, count uint64) error {
	err := k.EpochId.Set(ctx, count)
	if err != nil {
		return err
	}

	return nil
}

// AppendEpoch appends a epoch in the store with a new id and update the count
func (k Keeper) AppendEpoch(ctx context.Context, epoch types.Epoch) (uint64, error) {
	// Create the epoch
	id, err := k.EpochId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	epoch.Id = id

	err = k.SetEpoch(ctx, epoch)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetEpoch set a specific epoch in the store
func (k Keeper) SetEpoch(ctx context.Context, epoch types.Epoch) error {
	err := k.Epochs.Set(ctx, epoch.Id, epoch)
	if err != nil {
		return err
	}

	return nil
}

// GetEpoch returns a epoch from its id
func (k Keeper) GetEpoch(ctx context.Context, id uint64) (val types.Epoch, found bool, err error) {
	has, err := k.Epochs.Has(ctx, id)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Epochs.Get(ctx, id)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveEpoch removes a epoch from the store
func (k Keeper) RemoveEpoch(ctx context.Context, id uint64) error {
	err := k.Epochs.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllEpoch returns all epoch
func (k Keeper) GetAllEpoch(ctx context.Context) (list []types.Epoch, err error) {
	err = k.Epochs.Walk(
		ctx,
		nil,
		func(key uint64, value types.Epoch) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k Keeper) GetLastEpoch(ctx context.Context) (epoch types.Epoch, found bool, err error) {
	has, err := k.Epochs.Has(ctx, 0)
	if err != nil {
		return epoch, false, err
	}

	if !has {
		return epoch, false, nil
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
		return epoch, false, err
	}

	return epoch, true, nil
}
