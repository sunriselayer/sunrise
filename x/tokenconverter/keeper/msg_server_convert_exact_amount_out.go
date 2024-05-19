package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func (k msgServer) ConvertExactAmountOut(goCtx context.Context, msg *types.MsgConvertExactAmountOut) (*types.MsgConvertExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgConvertExactAmountOutResponse{}, nil
}
