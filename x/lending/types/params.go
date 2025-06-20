package types

import (
	"cosmossdk.io/math"
)

// NewParams creates a new Params instance.
func NewParams(ltvRatio, liquidationThreshold, baseInterestRate math.LegacyDec) Params {
	return Params{
		LtvRatio:             ltvRatio,
		LiquidationThreshold: liquidationThreshold,
		BaseInterestRate:     baseInterestRate,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		math.LegacyNewDecWithPrec(80, 2), // 0.80 = 80% LTV
		math.LegacyNewDecWithPrec(85, 2), // 0.85 = 85% liquidation threshold
		math.LegacyNewDecWithPrec(5, 2),  // 0.05 = 5% base interest rate
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.LtvRatio.IsNegative() || p.LtvRatio.GT(math.LegacyOneDec()) {
		return ErrInvalidLtvRatio
	}

	if p.LiquidationThreshold.IsNegative() || p.LiquidationThreshold.GT(math.LegacyOneDec()) {
		return ErrInvalidLiquidationThreshold
	}

	if p.LiquidationThreshold.LTE(p.LtvRatio) {
		return ErrInvalidLiquidationThreshold
	}

	if p.BaseInterestRate.IsNegative() {
		return ErrInvalidInterestRate
	}

	return nil
}
