package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k Keeper) swapRoute(
	ctx sdk.Context,
	sender sdk.AccAddress,
	route types.Route,
	amountExact math.Int,
	isExactAmountIn bool,
) (amountResult math.Int, err error) {
	switch strategy := route.Strategy.(type) {
	case *types.Route_Pool:
		amountResult, err = k.swapRoutePool(
			ctx,
			sender,
			strategy.Pool.PoolId,
			route.DenomIn,
			route.DenomOut,
			amountExact,
			isExactAmountIn,
		)
		if err != nil {
			return math.Int{}, err
		}

	case *types.Route_Series:
		for i := range strategy.Series.Routes {
			var r types.Route
			if isExactAmountIn {
				r = strategy.Series.Routes[i]
			} else {
				r = strategy.Series.Routes[len(strategy.Series.Routes)-1-i]
			}
			amountResult, err = k.swapRoute(ctx, sender, r, amountExact, isExactAmountIn)
			if err != nil {
				return math.Int{}, err
			}
			amountExact = amountResult
		}
		// No needs to do
		// amountResult = amountExact

	case *types.Route_Parallel:
		// Calculate the sum of the weights
		weightSum := math.LegacyZeroDec()
		for _, w := range strategy.Parallel.Weights {
			weightSum.AddMut(w)
		}

		// Calculate the amount of input for each route
		amountsExact := make([]math.Int, len(strategy.Parallel.Routes))
		amountsExactSum := math.ZeroInt()
		length := len(strategy.Parallel.Weights)

		for i, w := range strategy.Parallel.Weights[:length-1] {
			amountsExact[i] = w.MulInt(amountExact).Quo(weightSum).TruncateInt()
		}
		// For avoiding rounding errors
		amountsExact[length-1] = amountExact.Sub(amountsExactSum)

		// Execute the swaps
		amountsResultSum := math.ZeroInt()

		for i, r := range strategy.Parallel.Routes {
			aResult, err := k.swapRoute(ctx, sender, r, amountsExact[i], isExactAmountIn)
			if err != nil {
				return math.Int{}, err
			}
			amountsResultSum = amountsResultSum.Add(aResult)
		}
		amountResult = amountsResultSum
	}
	return amountResult, nil
}

func (k Keeper) swapRoutePool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	denomIn string,
	denomOut string,
	amountExact math.Int,
	isExactAmountIn bool,
) (amountResult math.Int, err error) {
	pool, found := k.liquidityPoolKeeper.GetPool(ctx, poolId)
	if !found {
		return math.Int{}, lptypes.ErrPoolNotFound
	}

	// No needs to validate the denom,
	// as liquiditypool side is responsible for ensuring the denom exists in the pool.

	if isExactAmountIn {
		tokenIn := sdk.NewCoin(denomIn, amountExact)

		amountResult, err = k.liquidityPoolKeeper.SwapExactAmountIn(
			ctx,
			sender,
			pool,
			tokenIn,
			denomOut,
		)
		if err != nil {
			return math.Int{}, err
		}
	} else {
		tokenOut := sdk.NewCoin(denomOut, amountExact)

		amountResult, err = k.liquidityPoolKeeper.SwapExactAmountOut(
			ctx,
			sender,
			pool,
			tokenOut,
			denomIn,
		)
		if err != nil {
			return math.Int{}, err
		}
	}

	return amountResult, nil
}
