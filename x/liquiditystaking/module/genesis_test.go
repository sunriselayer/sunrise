package liquiditystaking_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	liquiditystaking "github.com/sunriselayer/sunrise/x/liquiditystaking/module"
	"github.com/sunriselayer/sunrise/x/liquiditystaking/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LiquiditystakingKeeper(t)
	liquiditystaking.InitGenesis(ctx, k, genesisState)
	got := liquiditystaking.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
