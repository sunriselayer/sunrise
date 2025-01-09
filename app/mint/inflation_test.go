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
	genesis := ctx.BlockTime().Add(-366 * 24 * time.Hour)

	t.Run("normal inflation calculation", func(t *testing.T) {
		provision := mint.CalculateAnnualProvision(
			ctx,
			math.LegacyNewDecWithPrec(10, 2), // 10% initial inflation
			math.LegacyNewDecWithPrec(2, 2),  // 2% min inflation
			math.LegacyNewDecWithPrec(5, 2),  // 5% disinflation
			math.NewInt(1000000),
			genesis,
			math.NewInt(100000),
		)
		// inflation rate cap = 0.1 * (1 - 0.05)^1 = 0.095 (> 0.02)
		// next supply = (1 + 0.095) * 100000 = 109500
		// provision = 109500 - 100000 = 9500
		require.Equal(t, "9500", provision.String())
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
		// supply cap - total supply = 105000 - 100000 = 5000
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
