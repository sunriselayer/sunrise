package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func createNEpoch(keeper keeper.Keeper, ctx context.Context, n int) []types.Epoch {
	items := make([]types.Epoch, n)
	for i := range items {
		items[i].Id = keeper.AppendEpoch(ctx, items[i])
	}
	return items
}

func TestEpochGet(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	for _, item := range items {
		got, found := f.keeper.GetEpoch(f.ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestEpochRemove(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	for _, item := range items {
		f.keeper.RemoveEpoch(f.ctx, item.Id)
		_, found := f.keeper.GetEpoch(f.ctx, item.Id)
		require.False(t, found)
	}
}

func TestEpochGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(f.keeper.GetAllEpoch(f.ctx)),
	)
}

func TestEpochCount(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, f.keeper.GetEpochCount(f.ctx))
}

func TestGetLastEpoch(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	got, found := f.keeper.GetLastEpoch(f.ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&items[len(items)-1]),
		nullify.Fill(&got),
	)
}
