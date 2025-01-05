package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	if p.EpochBlocks <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch blocks must be positive")
	}

	if p.StakingRewardRatio.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "staking reward ratio must not be negative")
	}
	if p.StakingRewardRatio.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "staking reward ratio must be less than 1")
	}

	return nil
}
