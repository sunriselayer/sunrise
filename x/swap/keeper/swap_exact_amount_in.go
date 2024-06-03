package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) CalculateResultExactAmountIn(
	ctx sdk.Context,
	route types.Route,
	amountIn math.Int,
) (amountOut math.Int, err error) {
	panic("TODO")
}

func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	interfaceProvider string,
	route types.Route,
	amountIn math.Int,
	minAmountOut math.Int,
) (result types.RouteResult, interfaceFee math.Int, err error) {
	result, err = k.swapRouteExactAmountIn(ctx, sender, route, amountIn)

	if err != nil {
		return result, interfaceFee, err
	}

	if interfaceProvider != "" {
		addr, err := sdk.AccAddressFromBech32(interfaceProvider)
		if err != nil {
			return result, interfaceFee, err
		}

		params := k.GetParams(ctx)
		// amountOut = totalAmountOut - interfaceFee
		//           = totalAmountOut * (1 - interfaceFeeRate)
		totalAmountOut := result.TokenOut.Amount
		amountOut := math.LegacyNewDecFromInt(totalAmountOut).Mul(math.LegacyOneDec().Sub(params.InterfaceFeeRate)).TruncateInt()
		interfaceFee = totalAmountOut.Sub(amountOut)

		if result.TokenOut.Amount.LT(minAmountOut) {
			return result, interfaceFee, fmt.Errorf("TODO")
		}

		// TODO: Deduct interface fee
		_ = addr
	} else {
		interfaceFee = math.ZeroInt()
	}

	return result, interfaceFee, nil
}

func generateResultExactAmountIn(denomIn, denomOut string, amountExact, amountResult math.Int) (tokenIn sdk.Coin, tokenOut sdk.Coin) {
	return sdk.NewCoin(denomIn, amountExact), sdk.NewCoin(denomOut, amountResult)
}

func (k Keeper) swapRouteExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	route types.Route,
	amountIn math.Int,
) (result types.RouteResult, err error) {
	_, result, err = route.InspectRoute(
		amountIn,
		func(denomIn string, denomOut string, pool types.RoutePool, amountExact math.Int) (math.Int, error) {
			return k.swapRoutePoolExactAmountIn(ctx, sender, pool.PoolId, denomIn, denomOut, amountExact)
		},
		generateResultExactAmountIn,
		false,
	)

	return result, err
}

func (k Keeper) swapRoutePoolExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	denomIn string,
	denomOut string,
	amountIn math.Int,
) (amountOut math.Int, err error) {
	pool, found := k.liquidityPoolKeeper.GetPool(ctx, poolId)
	if !found {
		return math.Int{}, lptypes.ErrPoolNotFound
	}

	// No needs to validate the denom,
	// as liquiditypool side is responsible for ensuring the denom exists in the pool.
	tokenIn := sdk.NewCoin(denomIn, amountIn)

	amountOut, err = k.liquidityPoolKeeper.SwapExactAmountIn(
		ctx,
		sender,
		pool,
		tokenIn,
		denomOut,
	)
	if err != nil {
		return math.Int{}, err
	}

	return amountOut, nil
}
