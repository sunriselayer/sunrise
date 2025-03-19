package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/liquidstaking/types"
)

func (k msgServer) LiquidStake(ctx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgLiquidStakeResponse{}, nil
}
