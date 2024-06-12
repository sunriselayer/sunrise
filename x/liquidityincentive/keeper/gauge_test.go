package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNGauge(keeper keeper.Keeper, ctx context.Context, n int) []types.Gauge {
	items := make([]types.Gauge, n)
	for i := range items {
		items[i].PreviousEpochId = uint64(i)
		items[i].PoolId = uint64(i)

		keeper.SetGauge(ctx, items[i])
	}
	return items
}

func TestGaugeGet(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNGauge(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGauge(ctx, item.PreviousEpochId, item.PoolId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestGaugeRemove(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNGauge(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGauge(ctx, item.PreviousEpochId, item.PoolId)
		_, found := keeper.GetGauge(ctx,
			item.PreviousEpochId,
			item.PoolId,
		)
		require.False(t, found)
	}
}

func TestGaugeGetAll(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNGauge(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllGauges(ctx)),
	)
}
