package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	amountIn, err := k.Keeper.SwapExactAmountOut(ctx, sender, msg.InterfaceProvider, msg.Route, msg.MaxAmountIn, msg.AmountOut)
	if err != nil {
		return nil, err
	}

	return &types.MsgSwapExactAmountOutResponse{
		AmountIn: amountIn,
	}, nil
}
