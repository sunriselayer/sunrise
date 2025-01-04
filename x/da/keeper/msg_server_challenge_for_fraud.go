package keeper

import (
	"context"

	"sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) ChallengeForFraud(ctx context.Context, msg *types.MsgChallengeForFraud) (*types.MsgChallengeForFraudResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgChallengeForFraudResponse{}, nil
}
