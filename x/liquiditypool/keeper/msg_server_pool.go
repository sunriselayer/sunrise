package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pool = types.Pool{
		Sender:    msg.Sender,
		LowerTick: msg.LowerTick,
		UpperTick: msg.UpperTick,
	}

	id := k.AppendPool(
		ctx,
		pool,
	)

	return &types.MsgCreatePoolResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdatePool(goCtx context.Context, msg *types.MsgUpdatePool) (*types.MsgUpdatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pool = types.Pool{
		Sender:    msg.Sender,
		Id:        msg.Id,
		LowerTick: msg.LowerTick,
		UpperTick: msg.UpperTick,
	}

	// Checks that the element exists
	val, found := k.GetPool(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != val.Sender {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPool(ctx, pool)

	return &types.MsgUpdatePoolResponse{}, nil
}

func (k msgServer) DeletePool(goCtx context.Context, msg *types.MsgDeletePool) (*types.MsgDeletePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPool(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != val.Sender {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePool(ctx, msg.Id)

	return &types.MsgDeletePoolResponse{}, nil
}
