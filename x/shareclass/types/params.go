package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
)

// NewParams creates a new Params instance.
func NewParams(rewardPeriod time.Duration, createValidatorGas uint64) Params {
	return Params{
		RewardPeriod:       rewardPeriod,
		CreateValidatorGas: createValidatorGas,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(time.Hour*8, 1000000)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.RewardPeriod <= 0 {
		return errorsmod.Wrap(ErrInvalidRewardPeriod, "reward period must be greater than 0")
	}

	if p.CreateValidatorGas <= 0 {
		return errorsmod.Wrap(ErrInvalidCreateValidatorGas, "create validator gas must be greater than 0")
	}

	return nil
}
