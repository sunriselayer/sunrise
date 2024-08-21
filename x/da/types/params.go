package types

import (
	"encoding/base64"
	"time"

	"cosmossdk.io/math"
	"github.com/consensys/gnark-crypto/ecc"
	groth16 "github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sunriselayer/sunrise/x/da/zkp"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
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
		VoteThreshold:       voteThreshold,
		SlashEpoch:          slashEpoch,
		EpochMaxFault:       epochMaxFault,
		SlashFraction:       slashFraction,
		ReplicationFactor:   replicationFactor,
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

// DefaultParams returns a default set of parameters
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

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
