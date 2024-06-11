package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) VoteGauge(goCtx context.Context, msg *types.MsgVoteGauge) (*types.MsgVoteGaugeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgVoteGaugeResponse{}, nil
}
