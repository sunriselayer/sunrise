package keeper

import (
	"context"

    "github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


func (k msgServer) CreateTwap(goCtx context.Context,  msg *types.MsgCreateTwap) (*types.MsgCreateTwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value already exists
    _, isFound := k.GetTwap(
        ctx,
        msg.Index,
        )
    if isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
    }

    var twap = types.Twap{
        Creator: msg.Creator,
        Index: msg.Index,
        
    }

   k.SetTwap(
   		ctx,
   		twap,
   	)
	return &types.MsgCreateTwapResponse{}, nil
}

func (k msgServer) UpdateTwap(goCtx context.Context,  msg *types.MsgUpdateTwap) (*types.MsgUpdateTwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetTwap(
        ctx,
        msg.Index,
    )
    if !isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg creator is the same as the current owner
    if msg.Creator != valFound.Creator {
        return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    var twap = types.Twap{
		Creator: msg.Creator,
		Index: msg.Index,
        
	}

	k.SetTwap(ctx, twap)

	return &types.MsgUpdateTwapResponse{}, nil
}

func (k msgServer) DeleteTwap(goCtx context.Context,  msg *types.MsgDeleteTwap) (*types.MsgDeleteTwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetTwap(
        ctx,
        msg.Index,
    )
    if !isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg creator is the same as the current owner
    if msg.Creator != valFound.Creator {
        return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

	k.RemoveTwap(
	    ctx,
	msg.Index,
    )

	return &types.MsgDeleteTwapResponse{}, nil
}
