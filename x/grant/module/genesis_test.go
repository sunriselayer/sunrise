package grant_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "sunrise/testutil/keeper"
	"sunrise/testutil/nullify"
	"sunrise/x/grant/module"
	"sunrise/x/grant/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.GrantKeeper(t)
	grant.InitGenesis(ctx, k, genesisState)
	got := grant.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
