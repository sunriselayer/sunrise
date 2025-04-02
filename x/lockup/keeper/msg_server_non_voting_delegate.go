package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) NonVotingDelegate(ctx context.Context, msg *types.MsgNonVotingDelegate) (*types.MsgNonVotingDelegateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	// TODO: Handle the message

	return &types.MsgNonVotingDelegateResponse{}, nil
}
