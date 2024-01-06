package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

func (k msgServer) BurnDerivative(goCtx context.Context, msg *types.MsgBurnDerivative) (*types.MsgBurnDerivativeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}

	sharesReceived, err := k.Keeper.BurnDerivative(ctx, sender, validator, msg.Amount)
	if err != nil {
		return nil, err
	}

	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	// 	),
	// )

	return &types.MsgBurnDerivativeResponse{
		Received: sharesReceived,
	}, nil
}
