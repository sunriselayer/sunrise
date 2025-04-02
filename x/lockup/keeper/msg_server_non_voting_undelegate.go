package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func (k msgServer) NonVotingUndelegate(ctx context.Context, msg *types.MsgNonVotingUndelegate) (*types.MsgNonVotingUndelegateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Owner); err != nil {
		return nil, errorsmod.Wrap(err, "invalid owner address")
	}

	// TODO: Handle the message

	return &types.MsgNonVotingUndelegateResponse{}, nil
}
