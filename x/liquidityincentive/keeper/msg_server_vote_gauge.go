package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) VoteGauge(ctx context.Context, msg *types.MsgVoteGauge) (*types.MsgVoteGaugeResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgVoteGaugeResponse{}, nil
}
