package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iden3/go-iden3-crypto/poseidon"

	"github.com/sunriselayer/sunrise/x/da/das/kzg"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) RespondToChallenge(ctx context.Context, msg *types.MsgRespondToChallenge) (*types.MsgRespondToChallengeResponse, error) {
	_, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	_, found, err := k.GetChallenge(ctx, msg.ShardsMerkleRoot, msg.ShardIndex, msg.EvaluationPointIndex)
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

	hash := func(data []byte) [32]byte {
		return [32]byte(poseidon.Sum(data))
	}

	inclusion := types.InclusionProof(hash(msg.ShardsMerkleRoot), merklePath, msg.KzgCommitment, hash)
	if !inclusion {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify kzg commitment inclusion proof")
	}

	err = kzg.KzgVerify(msg.KzgCommitment, msg.KzgOpeningProof, int(msg.EvaluationPointIndex))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to verify kzg commitment")
	}

	return &types.MsgRespondToChallengeResponse{}, nil
}
