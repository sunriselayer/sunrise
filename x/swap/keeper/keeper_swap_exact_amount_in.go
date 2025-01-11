package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"

	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) calculateInterfaceFeeExactAmountIn(
	ctx sdk.Context,
	hasInterfaceFee bool,
	amountOutGross math.Int,
) (amountOutNet math.Int, interfaceFee math.Int) {
	if !hasInterfaceFee {
		return amountOutGross, math.ZeroInt()
	}

	//TODO: error handling
	params, _ := k.Params.Get(ctx)
	interfaceFeeRate := math.LegacyMustNewDecFromStr(params.InterfaceFeeRate) // TODO: remove with math.Dec
	// $ amountOutNet = amountOutGross - interfaceFee $
	//                = amountOutGross * (1 - interfaceFeeRate) $
	amountOutNet = math.LegacyNewDecFromInt(amountOutGross).Mul(math.LegacyOneDec().Sub(interfaceFeeRate)).TruncateInt()
	interfaceFee = amountOutGross.Sub(amountOutNet)

	return amountOutNet, interfaceFee
}

func (k Keeper) CalculateResultExactAmountIn(
	ctx sdk.Context,
	hasInterfaceFee bool,
	route types.Route,
	amountIn math.Int,
) (result types.RouteResult, interfaceFee math.Int, err error) {
	result, err = k.calculateResultRouteExactAmountIn(ctx, route, amountIn)
	if err != nil {
		return result, interfaceFee, err
	}

	var (
		amountOutGross = result.TokenOut.Amount
	)

	_, interfaceFee = k.calculateInterfaceFeeExactAmountIn(ctx, hasInterfaceFee, amountOutGross)

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
	result, err = k.swapRouteExactAmountIn(ctx, sender, route, amountIn)
	if err != nil {
		return result, interfaceFee, err
	}

	var (
		hasInterfaceFee = interfaceProvider != ""
		amountOutNet    math.Int
		amountOutGross  = result.TokenOut.Amount
	)

	amountOutNet, interfaceFee = k.calculateInterfaceFeeExactAmountIn(ctx, hasInterfaceFee, amountOutGross)

	if amountOutNet.LT(minAmountOut) {
		return result, interfaceFee, types.ErrLowerThanMinOutAmount
	}

	if hasInterfaceFee {
		// Validated in ValidateBasic
		addr, err := sdk.AccAddressFromBech32(interfaceProvider)
		if err != nil {
			return result, interfaceFee, err
		}
		fee := sdk.NewCoin(result.TokenOut.Denom, interfaceFee)

		if fee.IsPositive() {
			if err := k.BankKeeper.SendCoins(ctx, sender, addr, sdk.NewCoins(fee)); err != nil {
				return result, interfaceFee, err
			}
		}
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
		true,
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
		true,
	)
	if err != nil {
		return math.Int{}, err
	}

	return amountOut, nil
}
