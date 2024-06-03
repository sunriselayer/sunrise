package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	amountOut, err := k.Keeper.SwapExactAmountIn(ctx, sender, msg.Route, msg.AmountIn, msg.MinAmountOut)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapExactAmountInResponse{
		AmountOut: amountOut,
	}, nil
}
