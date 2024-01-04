package grant_test

import (
	"testing"

	keepertest "github.com/sunrise-zone/sunrise-app/testutil/keeper"
	"github.com/sunrise-zone/sunrise-app/testutil/nullify"
	grant "github.com/sunrise-zone/sunrise-app/x/blobgrant/module"
	"github.com/sunrise-zone/sunrise-app/x/blobgrant/types"

	"github.com/stretchr/testify/require"
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
