package keeper

import (
	"context"

	"sunrise/x/liquidityincentive/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) CollectVoteRewards(ctx context.Context, msg *types.MsgCollectVoteRewards) (*types.MsgCollectVoteRewardsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgCollectVoteRewardsResponse{}, nil
}
