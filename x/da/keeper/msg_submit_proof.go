package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
)

func (k msgServer) SubmitProof(goCtx context.Context, msg *types.MsgSubmitProof) (*types.MsgSubmitProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	publishedData := k.GetPublishedData(ctx, msg.MetadataUri)
	if publishedData.Status != "challenge_for_fraud" {
		return nil, types.ErrDataNotInChallenge
	}

	// check proof period
	params := k.GetParams(ctx)
	if publishedData.ChallengeTimestamp.Add(params.ProofPeriod).Before(ctx.BlockTime()) {
		return nil, types.ErrProofPeriodIsOver
	}

	// TODO: check proof
	// save proof in the storage
	err := k.SetProof(ctx, types.Proof{
		MetadataUri: msg.MetadataUri,
		Sender:      msg.Sender,
		Indices:     msg.Indices,
		ShardHashes: msg.ShardHashes,
	})
	if err != nil {
		return nil, err
	}

	err = ctx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSubmitProofResponse{}, nil
}
