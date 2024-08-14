package keeper

import (
	"bytes"
	"context"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	groth16bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/frontend"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/da/types"
	"github.com/sunriselayer/sunrise/x/da/zkp"
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

	// check proof
	if msg.IsValidData {
		// TODO: check number of proofs (threshold)
		vk, err := zkp.UnmarshalVerifyingKey(params.ZkpVerifyingKey)
		if err != nil {
			return nil, err
		}

		// groth16: Prove & Verify
		for i, j := range msg.Indices {
			proof := &groth16bn254.Proof{}
			_, err := proof.ReadFrom(bytes.NewReader(msg.Proofs[i]))
			if err != nil {
				return nil, err
			}

			if len(publishedData.ShardDoubleHashes) <= int(j) {
				return nil, types.ErrProofIndiceOverflow
			}

			assignment := zkp.ValidityProofCircuit{
				ShardHash:       big.NewInt(1),
				ShardDoubleHash: publishedData.ShardDoubleHashes[j],
			}
			witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
			if err != nil {
				return nil, err
			}

			pubWitness, err := witness.Public()
			if err != nil {
				return nil, err
			}

			err = groth16.Verify(proof, vk, pubWitness)
			if err != nil {
				return nil, err
			}
		}
	}

	// save proof in the storage
	err := k.SetProof(ctx, types.Proof{
		MetadataUri: msg.MetadataUri,
		Sender:      msg.Sender,
		Indices:     msg.Indices,
		Proofs:      msg.Proofs,
		IsValidData: msg.IsValidData,
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
