package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"sunrise/x/blob/types"
)

func (k msgServer) PayForBlobs(goCtx context.Context, msg *types.MsgPayForBlobs) (*types.MsgPayForBlobsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgPayForBlobsResponse{}, nil
}
