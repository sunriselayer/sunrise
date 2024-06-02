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
	amountIn math.Int,
) (amountOut math.Int, err error) {
	switch strategy := route.Strategy.(type) {
	case *types.Route_Pool:
		tokenIn := sdk.NewCoin(route.DenomIn, amountIn)
		amountOut, err = k.swapRoutePool(ctx, sender, strategy.Pool.PoolId, tokenIn, route.DenomOut)
		if err != nil {
			return math.Int{}, err
		}
	case *types.Route_Series:
		for _, r := range strategy.Series.Routes {
			amountOut, err = k.swapRoute(ctx, sender, r, amountIn)
			if err != nil {
				return math.Int{}, err
			}
			amountIn = amountOut
		}
	case *types.Route_Parallel:
		// Calculate the sum of the weights
		weightSum := math.LegacyZeroDec()
		for _, w := range strategy.Parallel.Weights {
			weightSum.AddMut(w)
		}

		// Calculate the amount of input for each route
		amountsIn := make([]math.Int, len(strategy.Parallel.Routes))
		amountsInSum := math.ZeroInt()
		length := len(strategy.Parallel.Weights)

		for i, w := range strategy.Parallel.Weights[:length-1] {
			amountsIn[i] = w.MulInt(amountIn).Quo(weightSum).TruncateInt()
		}
		// For avoiding rounding errors
		amountsIn[length-1] = amountIn.Sub(amountsInSum)

		// Execute the swaps
		amountsOutSum := math.ZeroInt()

		for i, r := range strategy.Parallel.Routes {
			aOut, err := k.swapRoute(ctx, sender, r, amountsIn[i])
			if err != nil {
				return math.Int{}, err
			}
			amountsOutSum = amountsOutSum.Add(aOut)
		}
		amountOut = amountsOutSum
	}
	return amountOut, nil
}

func (k Keeper) swapRoutePool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	tokenIn sdk.Coin,
	denomOut string,
) (amountOut math.Int, err error) {
	pool, found := k.liquidityPoolKeeper.GetPool(ctx, poolId)
	if !found {
		return math.Int{}, lptypes.ErrPoolNotFound
	}

	// No needs to validate the denom,
	// as liquiditypool side is responsible for ensuring the denom exists in the pool.
	amountOut, err = k.liquidityPoolKeeper.SwapExactAmountIn(ctx, sender, pool, tokenIn, denomOut, math.ZeroInt(), pool.FeeRate)
	if err != nil {
		return math.Int{}, err
	}

	return amountOut, nil
}
