package keeper

import (
	"context"
)

func (k Keeper) FeeDenom(ctx context.Context) (string, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return "", err
	}

	return params.FeeDenom, nil
}
