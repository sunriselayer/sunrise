package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNGauge(keeper keeper.Keeper, ctx context.Context, n int) []types.Gauge {
	items := make([]types.Gauge, n)
	for i := range items {
		items[i].PreviousEpochId = 1
		items[i].PoolId = uint64(i)
		items[i].Count = math.OneInt()

		_ = keeper.SetGauge(ctx, items[i])
	}
	return items
}

func TestGaugeGet(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	for _, item := range items {
		rst, found, err := f.keeper.GetGauge(f.ctx, item.PreviousEpochId, item.PoolId)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestGaugeRemove(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	for _, item := range items {
		err := f.keeper.RemoveGauge(f.ctx, item.PreviousEpochId, item.PoolId)
		require.NoError(t, err)
		_, found, err := f.keeper.GetGauge(f.ctx,
			item.PreviousEpochId,
			item.PoolId,
		)
		require.NoError(t, err)
		require.False(t, found)
	}
}

func TestGaugeGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	gauges, err := f.keeper.GetAllGauges(f.ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(gauges),
	)
}

func TestGetAllGaugeByPreviousEpochId(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	gauges, err := f.keeper.GetAllGaugeByPreviousEpochId(f.ctx, 1)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(gauges),
	)
	gauges, err = f.keeper.GetAllGaugeByPreviousEpochId(f.ctx, 2)
	require.NoError(t, err)
	require.Len(t,
		gauges,
		0,
	)
}
