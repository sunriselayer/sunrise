package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func (k msgServer) CollectIncentiveRewards(goCtx context.Context, msg *types.MsgCollectIncentiveRewards) (*types.MsgCollectIncentiveRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCollectIncentiveRewardsResponse{}, nil
}
