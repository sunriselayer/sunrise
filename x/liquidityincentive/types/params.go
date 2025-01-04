package types

import (
	"cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams(epochBlocks int64, stakingRewardRatio math.LegacyDec) Params {
	return Params{
		EpochBlocks:        epochBlocks,
		StakingRewardRatio: stakingRewardRatio,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		5,                                // new epoch per 10 blocks
		math.LegacyNewDecWithPrec(50, 2), // 50% to staking
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
