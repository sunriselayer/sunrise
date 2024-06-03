package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) calculateInterfaceFeeExactAmountIn(
	ctx sdk.Context,
	hasInterfaceFee bool,
	totalAmountOut math.Int,
) (amountOut math.Int, interfaceFee math.Int) {
	if !hasInterfaceFee {
		return totalAmountOut, math.ZeroInt()
	}

	params := k.GetParams(ctx)
	// amountOut = totalAmountOut - interfaceFee
	//           = totalAmountOut * (1 - interfaceFeeRate)
	amountOut = math.LegacyNewDecFromInt(totalAmountOut).Mul(math.LegacyOneDec().Sub(params.InterfaceFeeRate)).TruncateInt()
	interfaceFee = totalAmountOut.Sub(amountOut)

	return amountOut, interfaceFee
}

func (k Keeper) CalculateResultExactAmountIn(
	ctx sdk.Context,
	hasInterfaceFee bool,
	route types.Route,
	amountIn math.Int,
) (result types.RouteResult, interfaceFee math.Int, err error) {
	var (
		totalAmountOut = result.TokenOut.Amount
	)

	_, interfaceFee = k.calculateInterfaceFeeExactAmountIn(ctx, hasInterfaceFee, totalAmountOut)

	result, err = k.calculateResultRouteExactAmountIn(ctx, route, amountIn)
	if err != nil {
		return result, interfaceFee, err
	}

	return result, interfaceFee, nil
}

func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	interfaceProvider string,
	route types.Route,
	amountIn math.Int,
	minAmountOut math.Int,
) (result types.RouteResult, interfaceFee math.Int, err error) {

	var (
		hasInterfaceFee = interfaceProvider != ""
		amountOut       math.Int
		totalAmountOut  = result.TokenOut.Amount
	)

	amountOut, interfaceFee = k.calculateInterfaceFeeExactAmountIn(ctx, hasInterfaceFee, totalAmountOut)

	result, err = k.swapRouteExactAmountIn(ctx, sender, route, amountIn)
	if err != nil {
		return result, interfaceFee, err
	}

	if hasInterfaceFee {
		// Validated in ValidateBasic
		addr := sdk.MustAccAddressFromBech32(interfaceProvider)

		// TODO: Deduct interface fee
		_ = addr
	}

	if amountOut.LT(minAmountOut) {
		return result, interfaceFee, fmt.Errorf("TODO")
	}

	return result, interfaceFee, nil
}

func generateResultExactAmountIn(denomIn, denomOut string, amountExact, amountResult math.Int) (tokenIn sdk.Coin, tokenOut sdk.Coin) {
	return sdk.NewCoin(denomIn, amountExact), sdk.NewCoin(denomOut, amountResult)
}

func (k Keeper) calculateResultRouteExactAmountIn(
	ctx sdk.Context,
	route types.Route,
	amountIn math.Int,
) (result types.RouteResult, err error) {
	_, result, err = route.InspectRoute(
		amountIn,
		func(denomIn string, denomOut string, pool types.RoutePool, amountExact math.Int) (math.Int, error) {
			return k.calculateResultRoutePoolExactAmountIn(ctx, pool.PoolId, denomIn, denomOut, amountExact)
		},
		generateResultExactAmountIn,
		false,
	)

	return result, err
}

func (k Keeper) calculateResultRoutePoolExactAmountIn(
	ctx sdk.Context,
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

	amountOut, err = k.liquidityPoolKeeper.CalculateResultExactAmountIn(
		ctx,
		pool,
		tokenIn,
		denomOut,
	)
	if err != nil {
		return math.Int{}, err
	}

	return amountOut, nil
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
