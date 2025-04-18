package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance.
func NewParams(rewardPeriod time.Duration, createValidatorFee sdk.Coin) Params {
	return Params{
		RewardPeriod:       rewardPeriod,
		CreateValidatorFee: createValidatorFee,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(time.Hour*8, sdk.NewCoin("fee", math.NewInt(1000000)))
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.RewardPeriod <= 0 {
		return errorsmod.Wrap(ErrInvalidRewardPeriod, "reward period must be greater than 0")
	}

	if !p.CreateValidatorFee.IsValid() {
		return errorsmod.Wrap(ErrInvalidCreateValidatorFee, "create validator fee must be valid")
	}

	return nil
}
