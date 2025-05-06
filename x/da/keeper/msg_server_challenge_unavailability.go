package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/da/das/kzg"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) ChallengeUnavailability(ctx context.Context, msg *types.MsgChallengeUnavailability) (*types.MsgChallengeUnavailabilityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	blob, found, err := k.GetBlobCommitment(ctx, msg.ShardsMerkleRoot)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get blob commitment")
	}
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "blob commitment not found")
	}

	if msg.ShardIndex >= blob.ShardCount {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "shard index out of range")
	}

	if msg.EvaluationPointIndex >= kzg.EvaluationPointCount {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "evaluation point out of range (0-31)")
	}

	_, found, err = k.GetChallenge(ctx, msg.ShardsMerkleRoot, msg.ShardIndex, msg.EvaluationPointIndex)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get challenge")
	}
	if found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge already exists")

	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.BlockTime().After(blob.Expiry) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "blob commitment has expired")
	}

	err = k.SetChallenge(ctx, types.Challenge{
		ShardsMerkleRoot:     msg.ShardsMerkleRoot,
		ShardIndex:           msg.ShardIndex,
		EvaluationPointIndex: msg.EvaluationPointIndex,
	})
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to set challenge")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get params")
	}

	sdkCtx.GasMeter().ConsumeGas(params.GasChallengeUnavailability, "ChallengeUnavailability")

	return &types.MsgChallengeUnavailabilityResponse{}, nil
}
