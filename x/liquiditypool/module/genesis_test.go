package liquiditypool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	liquiditypool "github.com/sunriselayer/sunrise/x/liquiditypool/module"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Pools: []types.Pool{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		PoolCount: 2,
		Positions: []types.Position{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		PositionCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, _, ctx := keepertest.LiquiditypoolKeeper(t)
	liquiditypool.InitGenesis(ctx, k, genesisState)
	got := liquiditypool.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.Pools, got.Pools)
	require.Equal(t, genesisState.PoolCount, got.PoolCount)
	require.ElementsMatch(t, genesisState.Positions, got.Positions)
	require.Equal(t, genesisState.PositionCount, got.PositionCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
