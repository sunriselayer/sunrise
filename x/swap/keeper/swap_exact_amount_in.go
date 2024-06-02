package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	route types.Route,
	amountIn math.Int,
	minAmountOut math.Int,
) (amountOut math.Int, err error) {
	amountOut, err = k.swapRoute(ctx, sender, route, amountIn)

	if err != nil {
		return math.Int{}, err
	}

	if amountOut.LT(minAmountOut) {
		return math.Int{}, fmt.Errorf("TODO")
	}

	return amountOut, nil
}

func (k Keeper) EstimateExactAmountIn(
	ctx sdk.Context,
	route types.Route,
	amountIn math.Int,
) (amountOut math.Int, err error) {
	panic("TODO")
}
