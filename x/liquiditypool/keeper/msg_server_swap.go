package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)
	tokenIn := msg.TokenIn

	for _, route := range msg.Routes {
		amounts := types.AmountsFromWeights(route.PoolWeights, tokenIn.Amount)

		var tokenOut sdk.Coin
		if tokenIn.Denom == route.BaseDenom {
			tokenOut = sdk.NewCoin(route.QuoteDenom, math.ZeroInt())
		} else {
			tokenOut = sdk.NewCoin(route.BaseDenom, math.ZeroInt())
		}

		for i := range route.PoolWeights {
			amountOut, err := k.SwapExactAmountInSinglePool(
				ctx,
				route.PoolWeights[i].PoolId,
				sdk.NewCoin(tokenIn.Denom, amounts[i]),
				tokenOut.Denom,
				address,
				false,
			)

			if err != nil {
				return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error swapping in pool: %s", err.Error()))
			}

			tokenOut.Amount = tokenOut.Amount.Add(*amountOut)
		}

		tokenIn = tokenOut
	}
	tokenOut := tokenIn

	// check slippage can be done only for final token
	if tokenOut.Amount.LT(msg.MinAmountOut) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slippage exceeded")
	}

	// TODO: set PriceFootprint

	return &types.MsgSwapExactAmountInResponse{
		AmountOut: tokenOut.Amount,
	}, nil
}

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)
	tokenOut := msg.TokenOut
	route := msg.Route

	amounts := types.AmountsFromWeights(route.PoolWeights, tokenOut.Amount)

	var tokenIn sdk.Coin
	if tokenOut.Denom == route.BaseDenom {
		tokenIn = sdk.NewCoin(route.QuoteDenom, math.ZeroInt())
	} else {
		tokenIn = sdk.NewCoin(route.BaseDenom, math.ZeroInt())
	}

	for i := range route.PoolWeights {
		amountIn, err := k.SwapExactAmountOutSinglePool(
			ctx,
			route.PoolWeights[i].PoolId,
			sdk.NewCoin(tokenOut.Denom, amounts[i]),
			tokenIn.Denom,
			address,
			false,
		)

		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("error swapping in pool: %s", err.Error()))
		}

		tokenIn.Amount = tokenIn.Amount.Add(*amountIn)

		if tokenIn.Amount.GT(msg.MaxAmountIn) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slippage exceeded")
		}
	}

	// TODO: set PriceFootprint

	return &types.MsgSwapExactAmountOutResponse{
		AmountIn: tokenIn.Amount,
	}, nil
}
