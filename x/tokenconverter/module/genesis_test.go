package tokenconverter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	tokenconverter "github.com/sunriselayer/sunrise/x/tokenconverter/module"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TokenconverterKeeper(t)
	tokenconverter.InitGenesis(ctx, k, genesisState)
	got := tokenconverter.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
