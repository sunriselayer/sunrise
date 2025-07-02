package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return types.Params{}, err
	}
	return params, nil
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	return k.Params.Set(ctx, params)
}

func (k Keeper) GetNonTransferableDenom(ctx context.Context) (string, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return "", err
	}
	return params.NonTransferableDenom, nil
}

func (k Keeper) GetTransferableDenom(ctx context.Context) (string, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return "", err
	}
	return params.TransferableDenom, nil
}
