package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestTickToPrice_PositiveBaseOffset(t *testing.T) {
	tickParams := TickParams{
		PriceRatio: math.LegacyNewDecWithPrec(10001, 4).String(), // 1.0001
		BaseOffset: math.LegacyNewDecWithPrec(5, 1).String(),     // 0.5
	}
	tickPriceMultiplied, err := TickToMultipliedPrice(0, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000049998750062496.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000150003749937502.249600000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000250018750312495.999824960000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000350043752187527.249424942496000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1001050498891338999.924254447082603488")

	tickPriceMultiplied, err = TickToMultipliedPrice(-1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999950003749687527.247275272472752725")
	tickPriceMultiplied, err = TickToMultipliedPrice(-2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999850018747812745.972678004672285496")
	tickPriceMultiplied, err = TickToMultipliedPrice(-3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999750043743438402.132464758196465850")
	tickPriceMultiplied, err = TickToMultipliedPrice(-10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999050498558872229.393417455985515247")
}

func TestTickToPrice_NegativeBaseOffset(t *testing.T) {
	tickParams := TickParams{
		PriceRatio: math.LegacyNewDecWithPrec(10001, 4).String(), // 1.0001
		BaseOffset: math.LegacyNewDecWithPrec(-5, 1).String(),    // 0.5
	}
	tickPriceMultiplied, err := TickToMultipliedPrice(0, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999950003750000000.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000049998750375000.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000150003750250037.500000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000250018750625062.503750000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000950403851266689.899928761250000000")

	tickPriceMultiplied, err = TickToMultipliedPrice(-1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999850018748125187.481251874812518748")
	tickPriceMultiplied, err = TickToMultipliedPrice(-2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999750043743750812.400011873625156233")
	tickPriceMultiplied, err = TickToMultipliedPrice(-3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999650078735877224.677544119213234909")
	tickPriceMultiplied, err = TickToMultipliedPrice(-10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "998950603498834537.607465180925054338")
}

func TestTickToPrice_ZeroBaseOffset(t *testing.T) {
	tickParams := TickParams{
		PriceRatio: math.LegacyNewDecWithPrec(10001, 4).String(), // 1.0001
		BaseOffset: math.LegacyZeroDec().String(),                // 0
	}
	tickPriceMultiplied, err := TickToMultipliedPrice(0, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000000000000000000.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000100000000000000.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000200010000000000.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1000300030001000000.000000000000000000")
	tickPriceMultiplied, err = TickToMultipliedPrice(10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "1001000450120021003.000000000000000000")

	tickPriceMultiplied, err = TickToMultipliedPrice(-1, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999900009999000099.990000999900009999")
	tickPriceMultiplied, err = TickToMultipliedPrice(-2, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999800029996000499.940006999200089990")
	tickPriceMultiplied, err = TickToMultipliedPrice(-3, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999700059990001499.790027996400449945")
	tickPriceMultiplied, err = TickToMultipliedPrice(-10, tickParams)
	require.NoError(t, err)
	require.Equal(t, tickPriceMultiplied.String(), "999000549780071479.506172441398977632")
}

func TestGetSqrtPriceFromQuoteBase(t *testing.T) {
	sqrtPrice, err := GetSqrtPriceFromQuoteBase(math.NewInt(10000), math.NewInt(100))
	require.NoError(t, err)
	require.Equal(t, sqrtPrice.String(), "10.000000000000000000")

	sqrtPrice, err = GetSqrtPriceFromQuoteBase(math.NewInt(400), math.NewInt(100))
	require.NoError(t, err)
	require.Equal(t, sqrtPrice.String(), "2.000000000000000000")

	sqrtPrice, err = GetSqrtPriceFromQuoteBase(math.NewInt(1), math.NewInt(100))
	require.NoError(t, err)
	require.Equal(t, sqrtPrice.String(), "0.100000000000000000")

	sqrtPrice, err = GetSqrtPriceFromQuoteBase(math.NewInt(1), Multiplier.TruncateInt())
	require.NoError(t, err)
	require.Equal(t, sqrtPrice.String(), "0.000000001000000000")

	sqrtPrice, err = GetSqrtPriceFromQuoteBase(math.NewInt(1), Multiplier.Mul(Multiplier).TruncateInt())
	require.NoError(t, err)
	require.Equal(t, sqrtPrice.String(), "0.000000000000000001")
}
