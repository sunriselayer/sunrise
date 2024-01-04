package sunrise_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunrise-zone/sunrise-app/testutil/keeper"
	"github.com/sunrise-zone/sunrise-app/testutil/nullify"
	sunrise "github.com/sunrise-zone/sunrise-app/x/sunrise/module"
	"github.com/sunrise-zone/sunrise-app/x/sunrise/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SunriseKeeper(t)
	sunrise.InitGenesis(ctx, k, genesisState)
	got := sunrise.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
