package liquiditypool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
	"github.com/sunriselayer/sunrise-app/testutil/nullify"
	liquiditypool "github.com/sunriselayer/sunrise-app/x/liquiditypool/module"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PairList: []types.Pair{
		{
			Index: "0",
},
		{
			Index: "1",
},
	},
	PoolList: []types.Pool{
		{
			Id: 0,
		},
		{
			Id: 1,
		},
	},
	PoolCount: 2,
	TwapList: []types.Twap{
		{
			Index: "0",
},
		{
			Index: "1",
},
	},
	// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LiquiditypoolKeeper(t)
	liquiditypool.InitGenesis(ctx, k, genesisState)
	got := liquiditypool.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PairList, got.PairList)
require.ElementsMatch(t, genesisState.PoolList, got.PoolList)
require.Equal(t, genesisState.PoolCount, got.PoolCount)
require.ElementsMatch(t, genesisState.TwapList, got.TwapList)
// this line is used by starport scaffolding # genesis/test/assert
}
