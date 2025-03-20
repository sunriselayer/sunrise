package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
)

// NewParams creates a new Params instance.
func NewParams(rewardPeriod time.Duration) Params {
	return Params{
		RewardPeriod: rewardPeriod,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(time.Hour * 8)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.RewardPeriod <= 0 {
		return errorsmod.Wrap(ErrInvalidRewardPeriod, "reward period must be greater than 0")
	}

	return nil
}
