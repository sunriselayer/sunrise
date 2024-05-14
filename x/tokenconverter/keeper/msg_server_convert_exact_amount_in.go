package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) ConvertExactAmountIn(goCtx context.Context, msg *types.MsgConvertExactAmountIn) (*types.MsgConvertExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgConvertExactAmountInResponse{}, nil
}
