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

		EpochList: []types.Epoch{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		EpochCount: 2,
		GaugeList: []types.Gauge{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LiquidityincentiveKeeper(t)
	liquidityincentive.InitGenesis(ctx, k, genesisState)
	got := liquidityincentive.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.EpochList, got.EpochList)
	require.Equal(t, genesisState.EpochCount, got.EpochCount)
	require.ElementsMatch(t, genesisState.GaugeList, got.GaugeList)
	// this line is used by starport scaffolding # genesis/test/assert
}
