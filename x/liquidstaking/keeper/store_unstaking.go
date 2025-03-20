package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

// GetUnstakingId get the total number of Unstaking
func (k Keeper) GetUnstakingId(ctx context.Context) (id uint64, err error) {
	id, err = k.UnstakingId.Peek(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetUnstakingId set the total number of Unstaking
func (k Keeper) SetUnstakingId(ctx context.Context, id uint64) error {
	err := k.UnstakingId.Set(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// AppendUnstaking appends a unstaking in the store with a new id and update the count
func (k Keeper) AppendUnstaking(ctx context.Context, unstaking types.Unstaking) (id uint64, err error) {
	// Create the unstaking
	id, err = k.UnstakingId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	err = k.SetUnstaking(ctx, unstaking, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetPool set a specific pool in the store
func (k Keeper) SetUnstaking(ctx context.Context, unstaking types.Unstaking, id uint64) error {
	err := k.Unstakings.Set(ctx, collections.Join(unstaking.CompletionTime.Unix(), id), unstaking)
	if err != nil {
		return err
	}
	return nil
}

// GetUnstaking returns a unstaking from its id
func (k Keeper) GetUnstaking(ctx context.Context, completionTime time.Time, id uint64) (val types.Unstaking, found bool, err error) {
	has, err := k.Unstakings.Has(ctx, collections.Join(completionTime.Unix(), id))
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Unstakings.Get(ctx, collections.Join(completionTime.Unix(), id))
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveUnstaking removes a unstaking from the store
func (k Keeper) RemoveUnstaking(ctx context.Context, completionTime time.Time, id uint64) error {
	err := k.Unstakings.Remove(ctx, collections.Join(completionTime.Unix(), id))
	if err != nil {
		return err
	}
	return nil
}

// GetAllUnstakings returns all unstaking
func (k Keeper) GetAllUnstakings(ctx context.Context) (list []types.Unstaking, err error) {
	err = k.Unstakings.Walk(
		ctx,
		nil,
		func(key collections.Pair[int64, uint64], value types.Unstaking) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return
}

func (k Keeper) IterateCompletedUnstakings(ctx context.Context, now time.Time, cb func(id uint64, value types.Unstaking) (stop bool, err error)) error {
	return k.Unstakings.Walk(ctx, nil,
		func(key collections.Pair[int64, uint64], value types.Unstaking) (bool, error) {
			if value.CompletionTime.After(now) {
				return true, nil
			}

			return cb(key.K2(), value)
		},
	)
}
