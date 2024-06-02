package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	interfaceProvider string,
	route types.Route,
	maxAmountIn math.Int,
	amountOut math.Int,
) (amountIn math.Int, err error) {
	return amountIn, nil
}

func (k Keeper) CalculateResultExactAmountOut(
	ctx sdk.Context,
	route types.Route,
	amountOut math.Int,
) (amountIn math.Int, err error) {
	_ = lptypes.ModuleName

	return amountIn, nil
}
