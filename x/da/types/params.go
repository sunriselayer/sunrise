package types

import (
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
	maxShardSize uint64,
	challengePeriod time.Duration,
	proofPeriod time.Duration,
	challengeCollateral sdk.Coins,
	zkpProvingKey groth16.ProvingKey,
	zkpVerifyingKey groth16.VerifyingKey,
) Params {
	zkpProvingKeyBz, err := zkp.MarshalProvingKey(zkpProvingKey)
	if err != nil {
		panic(err)
	}
	zkpVerifyingKeyBz, err := zkp.MarshalProvingKey(zkpVerifyingKey)
	if err != nil {
		panic(err)
	}
	return Params{
		VoteThreshold:       voteThreshold,
		SlashEpoch:          slashEpoch,
		EpochMaxFault:       epochMaxFault,
		SlashFraction:       slashFraction,
		ReplicationFactor:   replicationFactor,
		MinShardCount:       minShardCount,
		MaxShardSize:        maxShardSize,
		ChallengePeriod:     challengePeriod,
		ProofPeriod:         proofPeriod,
		ChallengeCollateral: challengeCollateral,
		ZkpProvingKey:       zkpProvingKeyBz,
		ZkpVerifyingKey:     zkpVerifyingKeyBz,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &zkp.ValidityProofCircuit{})
	if err != nil {
		panic(err)
	}

	provingKey, verifyingKey, err := groth16.Setup(ccs)
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
		1000000,       // 1MB
		time.Minute*6, // 6min,
		time.Minute*8, // 8min
		sdk.Coins(nil),
		provingKey,
		verifyingKey,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
