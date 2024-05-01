package blobstream_test

import (
	"testing"

	keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
	"github.com/sunriselayer/sunrise-app/testutil/nullify"
	stream "github.com/sunriselayer/sunrise-app/x/blobstream/module"
	"github.com/sunriselayer/sunrise-app/x/blobstream/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.StreamKeeper(t)
	stream.InitGenesis(ctx, k, genesisState)
	got := stream.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
