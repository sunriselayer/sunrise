package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetValidatorsPowerSnapshot(ctx context.Context, blockHeight int64) (data types.ValidatorsPowerSnapshot, found bool, err error) {
	has, err := k.ValidatorsPowerSnapshots.Has(ctx, blockHeight)
	if err != nil {
		return data, false, err
	}

	if !has {
		return data, false, nil
	}

	val, err := k.ValidatorsPowerSnapshots.Get(ctx, blockHeight)
	if err != nil {
		return data, false, err
	}

	return val, true, nil
}

// SetParams set the params
func (k Keeper) SetValidatorsPowerSnapshot(ctx context.Context, data types.ValidatorsPowerSnapshot) error {
	err := k.ValidatorsPowerSnapshots.Set(ctx, data.BlockHeight, data)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteValidatorsPowerSnapshot(ctx context.Context, blockHeight int64) error {
	err := k.ValidatorsPowerSnapshots.Remove(ctx, blockHeight)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllValidatorsPowerSnapshots(ctx context.Context) (list []types.ValidatorsPowerSnapshot, err error) {
	err = k.ValidatorsPowerSnapshots.Walk(
		ctx,
		nil,
		func(key int64, value types.ValidatorsPowerSnapshot) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
