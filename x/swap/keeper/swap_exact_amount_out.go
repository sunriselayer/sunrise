package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k Keeper) CalculateResultExactAmountOut(
	ctx sdk.Context,
	route types.Route,
	amountOut math.Int,
) (result types.RouteResult, err error) {
	_ = lptypes.ModuleName

	return result, nil
}

func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	interfaceProvider string,
	route types.Route,
	maxAmountIn math.Int,
	amountOut math.Int,
) (result types.RouteResult, interfaceFee math.Int, err error) {
	if interfaceProvider != "" {
		addr, err := sdk.AccAddressFromBech32(interfaceProvider)
		if err != nil {
			return result, interfaceFee, err
		}

		params := k.GetParams(ctx)
		// totalAmountOut = amountOut + interfaceFee
		//                = amountOut / (1 - interfaceFeeRate)
		totalAmountOut := math.LegacyNewDecFromInt(amountOut).Quo(math.LegacyOneDec().Sub(params.InterfaceFeeRate)).TruncateInt()
		interfaceFee = totalAmountOut.Sub(amountOut)

		result, err = k.swapRouteExactAmountOut(ctx, sender, route, totalAmountOut)

		if err != nil {
			return result, interfaceFee, err
		}

		// TODO: Deduct interface fee
		_ = addr
	} else {
		interfaceFee = math.ZeroInt()
	}

	if result.TokenIn.Amount.GT(maxAmountIn) {
		return result, interfaceFee, fmt.Errorf("TODO")
	}

	return result, interfaceFee, nil
}

func generateResultExactAmountOut(denomIn, denomOut string, amountExact, amountResult math.Int) (tokenIn sdk.Coin, tokenOut sdk.Coin) {
	return sdk.NewCoin(denomIn, amountResult), sdk.NewCoin(denomOut, amountExact)
}

func (k Keeper) swapRouteExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	route types.Route,
	amountOut math.Int,
) (result types.RouteResult, err error) {
	_, result, err = route.InspectRoute(
		amountOut,
		func(denomIn string, denomOut string, pool types.RoutePool, amountExact math.Int) (math.Int, error) {
			return k.swapRoutePoolExactAmountOut(ctx, sender, pool.PoolId, denomIn, denomOut, amountExact)
		},
		generateResultExactAmountOut,
		false,
	)

	return result, err
}

func (k Keeper) swapRoutePoolExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	denomIn string,
	denomOut string,
	amountOut math.Int,
) (amountIn math.Int, err error) {
	pool, found := k.liquidityPoolKeeper.GetPool(ctx, poolId)
	if !found {
		return math.Int{}, lptypes.ErrPoolNotFound
	}

	// No needs to validate the denom,
	// as liquiditypool side is responsible for ensuring the denom exists in the pool.
	tokenOut := sdk.NewCoin(denomOut, amountOut)

	amountIn, err = k.liquidityPoolKeeper.SwapExactAmountOut(
		ctx,
		sender,
		pool,
		tokenOut,
		denomIn,
	)
	if err != nil {
		return math.Int{}, err
	}

	return amountIn, nil
}
