package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) ConvertExactAmountIn(goCtx context.Context, msg *types.MsgConvertExactAmountIn) (*types.MsgConvertExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	amountOut := k.Keeper.CalculateAmountOutFeeToken(ctx, msg.AmountIn)
	if amountOut.LT(msg.MinAmountOut) {
		return nil, types.ErrInsufficientAmountOut
	}

	if err := k.Keeper.Convert(ctx, msg.AmountIn, amountOut, address); err != nil {
		return nil, err
	}

	return &types.MsgConvertExactAmountInResponse{
		AmountOut: amountOut,
	}, nil
}
