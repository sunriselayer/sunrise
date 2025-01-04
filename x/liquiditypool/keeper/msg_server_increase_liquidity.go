package keeper

import (
	"context"

	"sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) IncreaseLiquidity(ctx context.Context, msg *types.MsgIncreaseLiquidity) (*types.MsgIncreaseLiquidityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgIncreaseLiquidityResponse{}, nil
}
