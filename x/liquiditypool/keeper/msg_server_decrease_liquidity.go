package keeper

import (
	"context"

	"sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) DecreaseLiquidity(ctx context.Context, msg *types.MsgDecreaseLiquidity) (*types.MsgDecreaseLiquidityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgDecreaseLiquidityResponse{}, nil
}
