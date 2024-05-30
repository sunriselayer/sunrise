package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestTickToPrice_PositiveBaseOffset(t *testing.T) {
	tickParams := TickParams{
		PriceRatio: math.LegacyNewDecWithPrec(10001, 4), // 1.0001
		BaseOffset: math.LegacyNewDecWithPrec(5, 1),     // 0.5
	}
	tickPrice, err := TickToPrice(0, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000049998750062496")
	tickPrice, err = TickToPrice(1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000150003749937502")
	tickPrice, err = TickToPrice(2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000250018750312496")
	tickPrice, err = TickToPrice(3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000350043752187527")
	tickPrice, err = TickToPrice(10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.001050498891339000")

	tickPrice, err = TickToPrice(-1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999950003749687527")
	tickPrice, err = TickToPrice(-2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999850018747812746")
	tickPrice, err = TickToPrice(-3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999750043743438402")
	tickPrice, err = TickToPrice(-10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999050498558872230")
}

func TestTickToPrice_NegativeBaseOffset(t *testing.T) {
	tickParams := TickParams{
		PriceRatio: math.LegacyNewDecWithPrec(10001, 4), // 1.0001
		BaseOffset: math.LegacyNewDecWithPrec(-5, 1),    // 0.5
	}
	tickPrice, err := TickToPrice(0, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999950003750000000")
	tickPrice, err = TickToPrice(1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000049998750375000")
	tickPrice, err = TickToPrice(2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000150003750250038")
	tickPrice, err = TickToPrice(3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000250018750625063")
	tickPrice, err = TickToPrice(10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000950403851266690")

	tickPrice, err = TickToPrice(-1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999850018748125187")
	tickPrice, err = TickToPrice(-2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999750043743750812")
	tickPrice, err = TickToPrice(-3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999650078735877225")
	tickPrice, err = TickToPrice(-10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.998950603498834538")
}

func TestTickToPrice_ZeroBaseOffset(t *testing.T) {
	tickParams := TickParams{
		PriceRatio: math.LegacyNewDecWithPrec(10001, 4), // 1.0001
		BaseOffset: math.LegacyZeroDec(),                // 0
	}
	tickPrice, err := TickToPrice(0, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000000000000000000")
	tickPrice, err = TickToPrice(1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000100000000000000")
	tickPrice, err = TickToPrice(2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000200010000000000")
	tickPrice, err = TickToPrice(3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.000300030001000000")
	tickPrice, err = TickToPrice(10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "1.001000450120021003")

	tickPrice, err = TickToPrice(-1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999900009999000100")
	tickPrice, err = TickToPrice(-2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999800029996000500")
	tickPrice, err = TickToPrice(-3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999700059990001500")
	tickPrice, err = TickToPrice(-10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPrice.String(), "0.999000549780071480")
}
