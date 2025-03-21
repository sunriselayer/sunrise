package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

// GetUnbondingId get the total number of Unbonding
func (k Keeper) GetUnbondingId(ctx context.Context) (id uint64, err error) {
	id, err = k.UnbondingId.Peek(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetUnbondingId set the total number of Unbonding
func (k Keeper) SetUnbondingId(ctx context.Context, id uint64) error {
	err := k.UnbondingId.Set(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// AppendUnbonding appends a unbonding in the store with a new id and update the count
func (k Keeper) AppendUnbonding(ctx context.Context, unbonding types.Unbonding) (id uint64, err error) {
	// Create the unbonding
	id, err = k.UnbondingId.Next(ctx)
	if err != nil {
		return 0, err
	}

	// Set the ID of the appended value
	err = k.SetUnbonding(ctx, unbonding, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SetUnbonding set a specific unbonding in the store
func (k Keeper) SetUnbonding(ctx context.Context, unbonding types.Unbonding, id uint64) error {
	err := k.Unbondings.Set(ctx, collections.Join(unbonding.CompletionTime.Unix(), id), unbonding)
	if err != nil {
		return err
	}
	return nil
}

// GetUnbonding returns a unbonding from its id
func (k Keeper) GetUnbonding(ctx context.Context, completionTime time.Time, id uint64) (val types.Unbonding, found bool, err error) {
	has, err := k.Unbondings.Has(ctx, collections.Join(completionTime.Unix(), id))
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Unbondings.Get(ctx, collections.Join(completionTime.Unix(), id))
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveUnbonding removes a unbonding from the store
func (k Keeper) RemoveUnbonding(ctx context.Context, completionTime time.Time, id uint64) error {
	err := k.Unbondings.Remove(ctx, collections.Join(completionTime.Unix(), id))
	if err != nil {
		return err
	}
	return nil
}

// GetAllUnbondings returns all unbonding
func (k Keeper) GetAllUnbondings(ctx context.Context) (list []types.Unbonding, err error) {
	err = k.Unbondings.Walk(
		ctx,
		nil,
		func(key collections.Pair[int64, uint64], value types.Unbonding) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return
}

func (k Keeper) IterateCompletedUnbondings(ctx context.Context, now time.Time, cb func(id uint64, value types.Unbonding) (stop bool, err error)) error {
	return k.Unbondings.Walk(ctx, nil,
		func(key collections.Pair[int64, uint64], value types.Unbonding) (bool, error) {
			if value.CompletionTime.After(now) {
				return true, nil
			}

			return cb(key.K2(), value)
		},
	)
}
