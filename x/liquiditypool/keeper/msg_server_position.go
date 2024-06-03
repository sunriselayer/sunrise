package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CreatePosition(goCtx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var position = types.Position{
		Address: msg.Sender,
	}

	id := k.AppendPosition(ctx, position)

	return &types.MsgCreatePositionResponse{
		Id:          id,
		AmountBase:  math.ZeroInt(), // TODO:
		AmountQuote: math.ZeroInt(), // TODO:
		Liquidity:   position.Liquidity,
	}, nil
}

func (k msgServer) IncreaseLiquidity(goCtx context.Context, msg *types.MsgIncreaseLiquidity) (*types.MsgIncreaseLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var position = types.Position{
		Address: msg.Sender,
		Id:      msg.Id,
	}

	// Checks that the element exists
	val, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != val.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPosition(ctx, position)

	return &types.MsgIncreaseLiquidityResponse{}, nil
}

func (k msgServer) DecreaseLiquidity(goCtx context.Context, msg *types.MsgDecreaseLiquidity) (*types.MsgDecreaseLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	position, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != position.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePosition(ctx, msg.Id)

	return &types.MsgDecreaseLiquidityResponse{}, nil
}
