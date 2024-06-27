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

		Epochs: []types.Epoch{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		EpochCount: 2,
		Gauges: []types.Gauge{
			{
				PreviousEpochId: 0,
				PoolId:          0,
			},
			{
				PreviousEpochId: 1,
				PoolId:          1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, _, ctx := keepertest.LiquidityincentiveKeeper(t)
	liquidityincentive.InitGenesis(ctx, k, genesisState)
	got := liquidityincentive.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.Epochs, got.Epochs)
	require.Equal(t, genesisState.EpochCount, got.EpochCount)
	require.ElementsMatch(t, genesisState.Gauges, got.Gauges)
	// this line is used by starport scaffolding # genesis/test/assert
}
