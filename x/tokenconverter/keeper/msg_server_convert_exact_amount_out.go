package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) ConvertExactAmountOut(goCtx context.Context, msg *types.MsgConvertExactAmountOut) (*types.MsgConvertExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	amountIn := k.Keeper.CalculateAmountInGovToken(ctx, msg.AmountOut)
	if amountIn.GT(msg.MaxAmountIn) {
		return nil, types.ErrExceededAmountIn
	}

	if err := k.Keeper.Convert(ctx, amountIn, msg.AmountOut, address); err != nil {
		return nil, err
	}

	return &types.MsgConvertExactAmountOutResponse{}, nil
}
