package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) CollectVoteRewards(ctx context.Context, msg *types.MsgCollectVoteRewards) (*types.MsgCollectVoteRewardsResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgCollectVoteRewardsResponse{}, nil
}
