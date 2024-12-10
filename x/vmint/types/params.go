package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	bondDenom string,
	feeDenom string,
	inflationRateCapInitial math.LegacyDec,
	inflationRateCapMinimum math.LegacyDec,
	disinflationRate math.LegacyDec,
	supplyCap math.Int,
	genesis time.Time,
) Params {
	return Params{
		BondDenom:               bondDenom,
		FeeDenom:                feeDenom,
		InflationRateCapInitial: inflationRateCapInitial.String(),
		InflationRateCapMinimum: inflationRateCapMinimum.String(),
		DisinflationRate:        disinflationRate.String(),
		SupplyCap:               supplyCap,
		Genesis:                 genesis,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"stake",
		"fee",
		math.LegacyMustNewDecFromStr("0.1"),
		math.LegacyMustNewDecFromStr("0.02"),
		math.LegacyMustNewDecFromStr("0.08"),
		math.NewInt(1_000_000_000).Mul(math.NewInt(1_000_000)),
		time.Now(),
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	err := sdk.ValidateDenom(p.BondDenom)
	if err != nil {
		return err
	}

	err = sdk.ValidateDenom(p.FeeDenom)
	if err != nil {
		return err
	}

	_, err = math.LegacyNewDecFromStr(p.InflationRateCapInitial)
	if err != nil {
		return err
	}

	_, err = math.LegacyNewDecFromStr(p.InflationRateCapMinimum)
	if err != nil {
		return err
	}

	_, err = math.LegacyNewDecFromStr(p.DisinflationRate)
	if err != nil {
		return err
	}

	if !p.SupplyCap.IsPositive() {
		return errorsmod.Wrap(ErrInvalidParam, "supply cap must be positive")
	}

	return nil
}
