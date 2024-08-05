package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) SubmitProof(goCtx context.Context, req *types.MsgSubmitProof) (*types.MsgSubmitProofResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx
	return &types.MsgSubmitProofResponse{}, nil
}
