package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

func (k msgServer) LiquidUnstake(goCtx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgLiquidUnstakeResponse{}, nil
}
