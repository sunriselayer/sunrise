package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/sunrise-zone/sunrise-app/testutil/keeper"
	"github.com/sunrise-zone/sunrise-app/x/sunrise/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.SunriseKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
