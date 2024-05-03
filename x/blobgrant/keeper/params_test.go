package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/blobgrant/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.GrantKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
