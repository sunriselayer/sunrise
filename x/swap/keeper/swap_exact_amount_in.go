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
	interfaceProvider string,
	route types.Route,
	amountIn math.Int,
	minAmountOut math.Int,
) (amountOut math.Int, err error) {
	amountOut, err = k.swapRoute(ctx, sender, route, amountIn, true)

	if err != nil {
		return math.Int{}, err
	}

	// TODO: Deduct interface provider fee

	if amountOut.LT(minAmountOut) {
		return math.Int{}, fmt.Errorf("TODO")
	}

	return amountOut, nil
}

func (k Keeper) CalculateResultExactAmountIn(
	ctx sdk.Context,
	route types.Route,
	amountIn math.Int,
) (amountOut math.Int, err error) {
	panic("TODO")
}
