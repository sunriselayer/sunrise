package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/swap/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SwapExactAmountIn(ctx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSwapExactAmountInResponse{}, nil
}
