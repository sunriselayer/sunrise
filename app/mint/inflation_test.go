package mint_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/app/mint"
)

func TestCalculateAnnualProvision(t *testing.T) {
	ctx := sdk.Context{}
	genesis := time.Now().Add(-365 * 24 * time.Hour)

	t.Run("normal inflation calculation", func(t *testing.T) {
		provision := mint.CalculateAnnualProvision(
			ctx,
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.LegacyNewDecWithPrec(2, 2),  // 2%
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.NewInt(1000000),
			genesis,
			math.NewInt(100000),
		)
		require.Equal(t, "9000", provision.String())
	})

	t.Run("hits minimum inflation rate", func(t *testing.T) {
		genesis := time.Now().Add(-10 * 365 * 24 * time.Hour)
		provision := mint.CalculateAnnualProvision(
			ctx,
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.LegacyNewDecWithPrec(2, 2),  // 2%
			math.LegacyNewDecWithPrec(20, 2), // 20%
			math.NewInt(1000000),
			genesis,
			math.NewInt(100000),
		)
		require.Equal(t, "2000", provision.String())
	})

	t.Run("hits supply cap", func(t *testing.T) {
		provision := mint.CalculateAnnualProvision(
			ctx,
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.LegacyNewDecWithPrec(2, 2),  // 2%
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.NewInt(105000),
			genesis,
			math.NewInt(100000),
		)
		require.Equal(t, "5000", provision.String())
	})

	t.Run("zero disinflation rate", func(t *testing.T) {
		provision := mint.CalculateAnnualProvision(
			ctx,
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.LegacyNewDecWithPrec(2, 2),  // 2%
			math.LegacyZeroDec(),
			math.NewInt(1000000),
			genesis,
			math.NewInt(100000),
		)
		require.Equal(t, "10000", provision.String())
	})

	t.Run("prevents supply decrease", func(t *testing.T) {
		provision := mint.CalculateAnnualProvision(
			ctx,
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.LegacyNewDecWithPrec(2, 2),  // 2%
			math.LegacyNewDecWithPrec(10, 2), // 10%
			math.NewInt(90000),               // Cap below current supply
			genesis,
			math.NewInt(100000),
		)
		require.Equal(t, "0", provision.String())
	})
}
