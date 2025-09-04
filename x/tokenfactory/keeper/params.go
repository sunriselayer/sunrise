package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/tokenfactory/types"
)

// SetParams sets the total supply of the community pool.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	return k.Params.Set(ctx, params)
}

// GetParams returns the total supply of the community pool.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}
	return params
}
