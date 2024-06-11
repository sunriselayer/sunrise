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
		items[i].Index = strconv.Itoa(i)

		keeper.SetGauge(ctx, items[i])
	}
	return items
}

func TestGaugeGet(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNGauge(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGauge(ctx,
			item.Index,
		)
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
		keeper.RemoveGauge(ctx,
			item.Index,
		)
		_, found := keeper.GetGauge(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestGaugeGetAll(t *testing.T) {
	keeper, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNGauge(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllGauge(ctx)),
	)
}
