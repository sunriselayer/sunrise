package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) ChallengeForFraud(goCtx context.Context, req *types.MsgChallengeForFraud) (*types.MsgChallengeForFraudResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	return &types.MsgChallengeForFraudResponse{}, nil
}
