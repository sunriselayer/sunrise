package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/da/das/kzg"
	"github.com/sunriselayer/sunrise/x/da/das/merkle"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) RespondToChallenge(ctx context.Context, msg *types.MsgRespondToChallenge) (*types.MsgRespondToChallengeResponse, error) {
	_, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	challenge, found, err := k.GetChallenge(ctx, msg.ChallengeId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get challenge")
	}
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "challenge not found")
	}

	merklePath := make([][33]byte, len(msg.KzgCommitmentMerklePath))
	for i, path := range msg.KzgCommitmentMerklePath {
		if len(path) != 33 {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "kzg commitment merkle path must be 33 bytes (1 byte 0x00 or 0x01 and 32 bytes hash)")
		}
		merklePath[i] = [33]byte(path)
	}

	inclusion := merkle.InclusionProof(merkle.Hash(challenge.ShardsMerkleRoot), merklePath, [32]byte(msg.KzgCommitment))
	if !inclusion {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify kzg commitment inclusion proof")
	}

	err = kzg.KzgVerify(msg.KzgCommitment, msg.KzgOpeningProof, int(challenge.EvaluationPointIndex))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to verify kzg commitment")
	}

	return &types.MsgRespondToChallengeResponse{}, nil
}
