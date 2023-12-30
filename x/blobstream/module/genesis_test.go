package stream_test

import (
	"testing"

	keepertest "sunrise/testutil/keeper"
	"sunrise/testutil/nullify"
	stream "sunrise/x/blobstream/module"
	"sunrise/x/blobstream/types"

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
