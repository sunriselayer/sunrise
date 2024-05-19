package liquidityincentive_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	liquidityincentive "github.com/sunriselayer/sunrise/x/liquidityincentive/module"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LiquidityincentiveKeeper(t)
	liquidityincentive.InitGenesis(ctx, k, genesisState)
	got := liquidityincentive.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
