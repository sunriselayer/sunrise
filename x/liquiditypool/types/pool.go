package types

import (
	"fmt"

	"cosmossdk.io/math"
)

func LpTokenDenom(poolId uint64) string {
	return fmt.Sprintf("%s/%d", ModuleName, poolId)
}

func LpTokenValueInQuoteUnit(baseAmount math.Int, quoteAmount math.Int, price math.LegacyDec) math.LegacyDec {
	baseAmountQuotedValue := baseAmount.ToLegacyDec().Mul(price)
	quoteAmountQuotedValue := quoteAmount.ToLegacyDec()

	return baseAmountQuotedValue.Add(quoteAmountQuotedValue)
}

func PoolModuleName(poolId uint64) string {
	return fmt.Sprintf("%s/%d", ModuleName, poolId)
}

func PoolTreasuryModuleName(poolId uint64) string {
	return fmt.Sprintf("%s/%d/treasury", ModuleName, poolId)
}

func CalculateX(y math.Int, k math.LegacyDec, f_x string) (*math.Int, error) {
	return nil, nil
}

func CalculateDx(y math.Int, dy math.Int, k math.LegacyDec, f_x string) (*math.Int, error) {
	left, err := CalculateX(y.Add(dy), k, f_x)
	if err != nil {
		return nil, err
	}

	right, err := CalculateX(y, k, f_x)
	if err != nil {
		return nil, err
	}

	diff := left.Sub(*right)

	return &diff, nil
}

func CalculateY(x math.Int, k math.LegacyDec, f_y string) (*math.Int, error) {
	return nil, nil
}

func CalculateDy(x math.Int, dx math.Int, k math.LegacyDec, f_y string) (*math.Int, error) {
	left, err := CalculateY(x.Add(dx), k, f_y)
	if err != nil {
		return nil, err
	}

	right, err := CalculateY(x, k, f_y)
	if err != nil {
		return nil, err
	}

	diff := left.Sub(*right)

	return &diff, nil
}

func CalculateK(x math.Int, y math.Int, f_k string) (*math.LegacyDec, error) {
	return nil, nil
}
