package keeper

import (
	"context"
)

func (k Keeper) GetFeeDenom(ctx context.Context) (string, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return "", err
	}
	return params.FeeDenom, nil
}

func (k Keeper) GetBondDenom(ctx context.Context) (string, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return "", err
	}
	return params.BondDenom, nil
}
