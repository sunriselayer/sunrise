package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CreatePosition(goCtx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var position = types.Position{
		Sender: msg.Sender,
	}

	id := k.AppendPosition(
		ctx,
		position,
	)

	return &types.MsgCreatePositionResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdatePosition(goCtx context.Context, msg *types.MsgUpdatePosition) (*types.MsgUpdatePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var position = types.Position{
		Sender: msg.Sender,
		Id:     msg.Id,
	}

	// Checks that the element exists
	val, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != val.Sender {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPosition(ctx, position)

	return &types.MsgUpdatePositionResponse{}, nil
}

func (k msgServer) DeletePosition(goCtx context.Context, msg *types.MsgDeletePosition) (*types.MsgDeletePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != val.Sender {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePosition(ctx, msg.Id)

	return &types.MsgDeletePositionResponse{}, nil
}
