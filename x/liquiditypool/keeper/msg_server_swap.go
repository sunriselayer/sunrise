package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)

	tokensOut, err := k.SwapExactAmountInMultiRoute(ctx, msg.Routes, msg.TokenIn, false, &address)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "error swapping: %s", err.Error())
	}

	tokenOut := tokensOut[len(tokensOut)-1]
	// check slippage can be done only for final token
	if tokenOut.Amount.LT(msg.MinAmountOut) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slippage exceeded")
	}

	// TODO: set PriceFootprint

	return &types.MsgSwapExactAmountInResponse{
		TokensVia: tokensOut[:len(tokensOut)-1],
		TokenOut:  tokenOut,
	}, nil
}

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)

	tokenIn, err := k.SwapExactAmountOutMultiPool(ctx, msg.Route, msg.TokenOut, false, &address, &msg.MaxAmountIn)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "error swapping: %s", err.Error())
	}

	// TODO: set PriceFootprint

	return &types.MsgSwapExactAmountOutResponse{
		TokenIn: *tokenIn,
	}, nil
}
