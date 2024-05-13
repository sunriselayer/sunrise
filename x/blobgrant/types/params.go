package types

import (
	"time"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	gasPerLiquidity math.LegacyDec,
	expiryDuration time.Duration,
	grantTokenRefillThreshold math.Int,
	blockHeightDuration uint64,
) Params {
	return Params{
		GasPerLiquidity:           gasPerLiquidity,
		ExpiryDuration:            expiryDuration,
		GrantTokenRefillThreshold: grantTokenRefillThreshold,
		BlockHeightDuration:       blockHeightDuration,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		math.LegacyNewDecWithPrec(1, 3),
		time.Hour*24*30,
		math.NewInt(1000_000),
		10,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.GasPerLiquidity.IsNil() || !p.GasPerLiquidity.IsPositive() {
		return ErrInvalidGasPerLiquidity
	}

	if p.ExpiryDuration == 0 {
		return ErrInvalidExpiryDuration
	}

	if p.GrantTokenRefillThreshold.IsNil() || !p.GrantTokenRefillThreshold.IsPositive() {
		return ErrInvalidGrantTokenRefillThreshold
	}

	if p.BlockHeightDuration == 0 {
		return ErrInvalidBlockHeightDuration
	}

	return nil
}
