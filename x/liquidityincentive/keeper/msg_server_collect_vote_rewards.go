package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) CollectVoteRewards(goCtx context.Context, msg *types.MsgCollectVoteRewards) (*types.MsgCollectVoteRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCollectVoteRewards{
		Address: msg.Sender,
	}); err != nil {
		return nil, err
	}

	return &types.MsgCollectVoteRewardsResponse{}, nil
}
