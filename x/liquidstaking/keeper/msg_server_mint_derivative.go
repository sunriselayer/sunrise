package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

func (k msgServer) MintDerivative(goCtx context.Context, msg *types.MsgMintDerivative) (*types.MsgMintDerivativeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	mintedDerivative, err := k.Keeper.MintDerivative(ctx, sender, validator, msg.Amount)
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

	return &types.MsgMintDerivativeResponse{
		Received: mintedDerivative,
	}, nil
}
