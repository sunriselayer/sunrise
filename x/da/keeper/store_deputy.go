package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k Keeper) GetDeputy(ctx context.Context, validator sdk.ValAddress) (deputy types.Deputy, found bool, err error) {
	has, err := k.Deputies.Has(ctx, validator)
	if err != nil {
		return deputy, false, err
	}

	if !has {
		return deputy, false, nil
	}

	deputy, err = k.Deputies.Get(ctx, validator)
	if err != nil {
		return deputy, false, err
	}

	return deputy, true, nil
}

// SetProofDeputy set the proof deputy of the validator
func (k Keeper) SetDeputy(ctx context.Context, validator sdk.ValAddress, deputy string) error {
	err := k.Deputies.Set(ctx, validator, types.Deputy{
		Validator: validator.String(),
		Address:   deputy,
	})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) DeleteDeputy(ctx context.Context, validator sdk.ValAddress) error {
	err := k.Deputies.Remove(ctx, validator)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetAllDeputies(ctx context.Context) (list []types.Deputy, err error) {
	err = k.Deputies.Walk(
		ctx,
		nil,
		func(key []byte, value types.Deputy) (bool, error) {
			list = append(list, value)
			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
