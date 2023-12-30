package blob_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "sunrise/testutil/keeper"
	"sunrise/testutil/nullify"
	"sunrise/x/blob/module"
	"sunrise/x/blob/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BlobKeeper(t)
	blob.InitGenesis(ctx, k, genesisState)
	got := blob.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
