package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ClaimRewards(ctx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgClaimRewardsResponse{}, nil
}
