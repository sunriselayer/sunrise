package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) Convert(goCtx context.Context, msg *types.MsgConvert) (*types.MsgConvertResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	amount, err := k.Keeper.CalculateConversionAmount(ctx, msg.MinAmount, msg.MaxAmount)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.Convert(ctx, amount, address); err != nil {
		return nil, err
	}

	return &types.MsgConvertResponse{
		Amount: amount,
	}, nil
}
