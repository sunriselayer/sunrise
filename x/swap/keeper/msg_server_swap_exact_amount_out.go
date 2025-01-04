package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/swap/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SwapExactAmountOut(ctx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSwapExactAmountOutResponse{}, nil
}
