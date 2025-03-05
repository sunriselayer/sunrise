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

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/da/zkp"
)

// NewParams creates a new Params instance.
func NewParams(
	challengeThreshold math.LegacyDec,
	replicationFactor math.LegacyDec,
	slashEpoch uint64,
	SlashFaultThreshold math.LegacyDec,
	slashFraction math.LegacyDec,
	challengePeriod time.Duration,
	proofPeriod time.Duration,
	rejectedRemovalPeriod time.Duration,
	publishDataCollateral sdk.Coins,
	submitInvalidityCollateral sdk.Coins,
	zkpVerifyingKey []byte,
	zkpProvingKey []byte,
	minShardCount uint64,
	maxShardCount uint64,
	maxShardSize uint64,
) Params {
	return Params{
		ChallengeThreshold:         challengeThreshold.String(),
		ReplicationFactor:          replicationFactor.String(),
		SlashEpoch:                 slashEpoch,
		SlashFaultThreshold:        SlashFaultThreshold.String(),
		SlashFraction:              slashFraction.String(),
		ChallengePeriod:            challengePeriod,
		ProofPeriod:                proofPeriod,
		RejectedRemovalPeriod:      rejectedRemovalPeriod,
		PublishDataCollateral:      publishDataCollateral,
		SubmitInvalidityCollateral: submitInvalidityCollateral,
		ZkpVerifyingKey:            zkpVerifyingKey,
		ZkpProvingKey:              zkpProvingKey,
		MinShardCount:              minShardCount,
		MaxShardCount:              maxShardCount,
		MaxShardSize:               maxShardSize,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	verifyingKey, err := base64.StdEncoding.DecodeString(DefaultVerifyingKeyBase64)
	if err != nil {
		panic(err)
	}

	provingKey, err := base64.StdEncoding.DecodeString(DefaultProvingKeyBase64)
	if err != nil {
		panic(err)
	}

	return NewParams(
		math.LegacyNewDecWithPrec(33, 2), // 33% challenge threshold
		math.LegacyNewDec(5),             // 5.0 replication factor
		120960,                           // 1 week(5sec/block) slash epoch
		math.LegacyNewDecWithPrec(5, 1),  // 50% slash fault threshold
		math.LegacyNewDecWithPrec(1, 3),  // 0.1% slash fraction
		time.Minute*4,                    // challenge 4min,
		time.Minute*10,                   // proof 10min
		time.Hour*24,                     // rejected remove 24h
		sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, math.NewInt(1_000_000_000))), // publish data collateral 1000RISE
		sdk.NewCoins(sdk.NewCoin(consts.FeeDenom, math.NewInt(100_000_000))),   // submit invalidity collateral 100RISE
		verifyingKey,
		provingKey,
		10,
		255,
		1000000, // 1MB
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

	zkpVerifyingKeyBz, err := zkp.MarshalVerifyingKey(verifyingKey)
	if err != nil {
		panic(err)
	}

	provingKeyBase64 := base64.StdEncoding.EncodeToString(zkpProvingKeyBz)
	verifyingKeyBase64 := base64.StdEncoding.EncodeToString(zkpVerifyingKeyBz)
	return provingKeyBase64, verifyingKeyBase64
}

// Validate validates the set of params.
func (p Params) Validate() error {
	challengeThreshold, err := math.LegacyNewDecFromStr(p.ChallengeThreshold)
	if err != nil {
		return err
	}
	if challengeThreshold.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge threshold must not be negative")
	}
	if challengeThreshold.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge threshold must be less than 1")
	}

	replicationFactor, err := math.LegacyNewDecFromStr(p.ReplicationFactor)
	if err != nil {
		return err
	}
	if !replicationFactor.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "replication factor must be negative")
	}

	if p.SlashEpoch == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slash epoch must be positive")
	}

	SlashFaultThreshold, err := math.LegacyNewDecFromStr(p.SlashFaultThreshold)
	if err != nil {
		return err
	}
	if SlashFaultThreshold.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slash fault threshold must not be negative")
	}
	if SlashFaultThreshold.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "slash fault threshold must be less than 1")
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

	if p.ChallengePeriod <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge period must be positive")
	}
	if p.ProofPeriod <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "proof period must be positive")
	}
	if p.RejectedRemovalPeriod <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "rejected removal period must be positive")
	}

	if !p.PublishDataCollateral.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "publish data collateral must be valid")
	}
	if !p.SubmitInvalidityCollateral.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "submit invalidity collateral must be valid")
	}

	if len(p.ZkpVerifyingKey) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "zkp verifying key must not be empty")
	}

	if len(p.ZkpProvingKey) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "zkp proving key must not be empty")
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

	return nil
}
