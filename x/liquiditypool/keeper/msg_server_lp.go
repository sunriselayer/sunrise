package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool) (*types.MsgJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	shareAmount, err := k.joinPool(ctx, msg.PoolId, msg.BaseToken, msg.QuoteToken, false, &msg.Sender, &msg.MinShareAmount)
	if err != nil {
		return nil, err
	}

	return &types.MsgJoinPoolResponse{
		ShareAmount: *shareAmount,
	}, nil
}

func (k msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokensOut, err := k.exitPool(ctx, msg.PoolId, msg.ShareAmount, false, &msg.Sender, &msg.MinAmountBase, &msg.MinAmountQuote)
	if err != nil {
		return nil, err
	}

	return &types.MsgExitPoolResponse{
		TokensOut: tokensOut,
	}, nil
}
