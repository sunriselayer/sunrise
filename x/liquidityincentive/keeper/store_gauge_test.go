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

		keeper.SetGauge(ctx, items[i])
	}
	return items
}

func TestGaugeGet(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	for _, item := range items {
		rst, found := f.keeper.GetGauge(f.ctx, item.PreviousEpochId, item.PoolId)
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
		f.keeper.RemoveGauge(f.ctx, item.PreviousEpochId, item.PoolId)
		_, found := f.keeper.GetGauge(f.ctx,
			item.PreviousEpochId,
			item.PoolId,
		)
		require.False(t, found)
	}
}

func TestGaugeGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(f.keeper.GetAllGauges(f.ctx)),
	)
}

func TestGetAllGaugeByPreviousEpochId(t *testing.T) {
	f := initFixture(t)
	items := createNGauge(f.keeper, f.ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(f.keeper.GetAllGaugeByPreviousEpochId(f.ctx, 1)),
	)
	require.Len(t,
		f.keeper.GetAllGaugeByPreviousEpochId(f.ctx, 2),
		0,
	)
}
