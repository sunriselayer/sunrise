package keeper

import (
	"bytes"
	"context"
	"math/big"

	"github.com/sunriselayer/sunrise/x/da/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	groth16bn254 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/frontend"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/zkp"
)

func (k msgServer) SubmitValidityProof(ctx context.Context, msg *types.MsgSubmitValidityProof) (*types.MsgSubmitValidityProofResponse, error) {
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	validator, err := k.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid validator address")
	}
	valAddr := sdk.ValAddress(validator)
	_, err = k.StakingKeeper.Validator(ctx, valAddr)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(sender, validator) {
		deputy, found := k.GetProofDeputy(ctx, validator)
		if !found {
			return nil, types.ErrDeputyNotFound
		}
		if !bytes.Equal(deputy, sender) {
			return nil, types.ErrInvalidDeputy
		}
	}

	// check number of proofs <> indices
	if len(msg.Indices) != len(msg.Proofs) {
		return nil, types.ErrIndicesAndProofsMismatch
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	publishedData, found := k.GetPublishedData(ctx, msg.MetadataUri)
	if !found {
		return nil, types.ErrDataNotFound
	}
	if publishedData.Status != types.Status_STATUS_CHALLENGING {
		return nil, types.ErrDataNotInChallenge
	}

	// check proof period
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	if publishedData.Timestamp.Add(params.ChallengePeriod).Add(params.ProofPeriod).Before(sdkCtx.BlockTime()) {
		return nil, types.ErrProofPeriodIsOver
	}

	// check proof
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
			return nil, types.ErrProofIndicesOverflow
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

	// save proof in the storage
	err = k.SetProof(ctx, types.Proof{
		MetadataUri: msg.MetadataUri,
		Sender:      sdk.AccAddress(validator).String(),
		Indices:     msg.Indices,
		Proofs:      msg.Proofs,
	})
	if err != nil {
		return nil, err
	}

	err = sdkCtx.EventManager().EmitTypedEvent(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitValidityProofResponse{}, nil
}
