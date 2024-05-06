package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)

	tokensVia, tokenOut, err := k.SwapExactAmountInMultiRoute(ctx, msg.Routes, msg.TokenIn, false, &address, &msg.MinAmountOut)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "error swapping: %s", err.Error())
	}

	return &types.MsgSwapExactAmountInResponse{
		TokensVia: tokensVia,
		TokenOut:  *tokenOut,
	}, nil
}

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address := sdk.MustAccAddressFromBech32(msg.Sender)

	tokenIn, err := k.SwapExactAmountOutMultiPool(ctx, msg.Route, msg.TokenOut, false, &address, &msg.MaxAmountIn)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "error swapping: %s", err.Error())
	}

	return &types.MsgSwapExactAmountOutResponse{
		TokenIn: *tokenIn,
	}, nil
}
