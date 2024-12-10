package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(epochBlocks int64, stakingRewardRatio math.LegacyDec) Params {
	return Params{
		EpochBlocks:        epochBlocks,
		StakingRewardRatio: stakingRewardRatio.String(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		10,                               // new epoch per 10 blocks
		math.LegacyNewDecWithPrec(50, 2), // 50% to staking
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.EpochBlocks <= 0 {
		return errorsmod.Wrap(ErrInvalidParam, "EpochBlocks must be positive")
	}

	_, err := math.LegacyNewDecFromStr(p.StakingRewardRatio)
	if err != nil {
		return err
	}

	return nil
}
