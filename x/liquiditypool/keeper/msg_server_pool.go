package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pool = types.Pool{
		Creator:    msg.Creator,
		BaseDenom:  msg.BaseDenom,
		QuoteDenom: msg.QuoteDenom,
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
		Creator:    msg.Creator,
		Id:         msg.Id,
		BaseDenom:  msg.BaseDenom,
		QuoteDenom: msg.QuoteDenom,
	}

	// Checks that the element exists
	val, found := k.GetPool(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPool(ctx, pool)

	return &types.MsgUpdatePoolResponse{}, nil
}
