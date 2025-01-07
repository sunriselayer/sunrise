package types

import (
	"encoding/base64"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	groth16 "github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/da/zkp"
)

// NewParams creates a new Params instance.
func NewParams(
	voteThreshold math.LegacyDec,
	slashEpoch uint64,
	epochMaxFault uint64,
	slashFraction math.LegacyDec,
	replicationFactor math.LegacyDec,
	minShardCount uint64,
	maxShardCount uint64,
	maxShardSize uint64,
	challengePeriod time.Duration,
	proofPeriod time.Duration,
	challengeCollateral sdk.Coins,
	zkpProvingKey []byte,
	zkpVerifyingKey []byte,
) Params {
	return Params{
		VoteThreshold:       voteThreshold.String(),
		SlashEpoch:          slashEpoch,
		EpochMaxFault:       epochMaxFault,
		SlashFraction:       slashFraction.String(),
		ReplicationFactor:   replicationFactor.String(),
		MinShardCount:       minShardCount,
		MaxShardCount:       maxShardCount,
		MaxShardSize:        maxShardSize,
		ChallengePeriod:     challengePeriod,
		ProofPeriod:         proofPeriod,
		ChallengeCollateral: challengeCollateral,
		ZkpProvingKey:       zkpProvingKey,
		ZkpVerifyingKey:     zkpVerifyingKey,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	provingKey, err := base64.StdEncoding.DecodeString(DefaultProvingKeyBase64)
	if err != nil {
		panic(err)
	}

	verifyingKey, err := base64.StdEncoding.DecodeString(DefaultVerifyingKeyBase64)
	if err != nil {
		panic(err)
	}

	return NewParams(
		math.LegacyNewDecWithPrec(67, 2), // 67%
		120960,                           // 1 week
		34560,                            // 2 days
		math.LegacyNewDecWithPrec(1, 3),  // 0.1%
		math.LegacyNewDec(5),             // 5.0
		10,
		255,
		1000000,       // 1MB
		time.Minute*6, // 6min,
		time.Minute*8, // 8min
		sdk.Coins(nil),
		provingKey,
		verifyingKey,
	)
}

func GenerateZkpKeys() (string, string) {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &zkp.ValidityProofCircuit{})
	if err != nil {
		panic(err)
	}

	provingKey, verifyingKey, err := groth16.Setup(ccs)
	if err != nil {
		panic(err)
	}

	zkpProvingKeyBz, err := zkp.MarshalProvingKey(provingKey)
	if err != nil {
		panic(err)
	}

	zkpVerifyingKeyBz, err := zkp.MarshalProvingKey(verifyingKey)
	if err != nil {
		panic(err)
	}

	provingKeyBase64 := base64.StdEncoding.EncodeToString(zkpProvingKeyBz)
	verifyingKeyBase64 := base64.StdEncoding.EncodeToString(zkpVerifyingKeyBz)
	return provingKeyBase64, verifyingKeyBase64
}

// Validate validates the set of params.
func (p Params) Validate() error {
	voteThreshold, err := math.LegacyNewDecFromStr(p.VoteThreshold)
	if err != nil {
		return err
	}
	if voteThreshold.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "vote threshold must not be negative")
	}
	if voteThreshold.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "vote threshold must be less than 1")
	}

	if p.SlashEpoch == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slash epoch must be positive")
	}

	if p.EpochMaxFault == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch max fault must be positive")
	}

	slashFraction, err := math.LegacyNewDecFromStr(p.SlashFraction)
	if err != nil {
		return err
	}
	if slashFraction.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slash fraction must not be negative")
	}
	if slashFraction.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slash fraction must be less than 1")
	}

	replicationFactor, err := math.LegacyNewDecFromStr(p.ReplicationFactor)
	if err != nil {
		return err
	}
	if replicationFactor.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "replication factor must not be negative")
	}
	if replicationFactor.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "replication factor must be less than 1")
	}

	if p.MinShardCount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "min shard count must be positive")
	}

	if p.MaxShardCount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max shard count must be positive")
	}

	if p.MaxShardCount < p.MinShardCount {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max shard count must be greater than or equal to min shard count")
	}

	if p.MaxShardSize == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max shard size must be positive")
	}

	if p.ChallengePeriod == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge period must be positive")
	}

	if p.ProofPeriod == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "proof period must be positive")
	}

	if !p.ChallengeCollateral.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge collateral must be valid")
	}

	if len(p.ZkpProvingKey) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "zkp proving key must not be empty")
	}

	if len(p.ZkpVerifyingKey) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "zkp verifying key must not be empty")
	}

	return nil
}
