package keeper

import (
	"context"

	"sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) SubmitProof(ctx context.Context, msg *types.MsgSubmitProof) (*types.MsgSubmitProofResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	// TODO: Handle the message

	return &types.MsgSubmitProofResponse{}, nil
}
