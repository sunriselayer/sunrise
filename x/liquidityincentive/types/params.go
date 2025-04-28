package types

import (
	"cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(epochBlocks int64, stakingRewardRatio math.LegacyDec, bribeClaimEpochs uint64) Params {
	return Params{
		EpochBlocks:        epochBlocks,
		StakingRewardRatio: stakingRewardRatio.String(),
		BribeClaimEpochs:   bribeClaimEpochs,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		7200,                             // new epoch per about a day
		math.LegacyNewDecWithPrec(50, 2), // 50% to staking
		5,                                // 5 epochs to claim bribes
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.EpochBlocks <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch blocks must be positive")
	}

	stakingRewardRatio, err := math.LegacyNewDecFromStr(p.StakingRewardRatio)
	if err != nil {
		return err
	}
	if stakingRewardRatio.IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "staking reward ratio must not be negative")
	}
	if stakingRewardRatio.GT(math.LegacyOneDec()) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "staking reward ratio must be less than 1")
	}

	if p.BribeClaimEpochs == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bribe claim epochs must be positive")
	}

	return nil
}
