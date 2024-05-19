package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k msgServer) CollectFees(goCtx context.Context, msg *types.MsgCollectFees) (*types.MsgCollectFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCollectFeesResponse{}, nil
}
