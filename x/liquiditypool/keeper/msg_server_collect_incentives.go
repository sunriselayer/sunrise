package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CollectIncentives(goCtx context.Context, msg *types.MsgCollectIncentives) (*types.MsgCollectIncentivesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCollectIncentivesResponse{}, nil
}
