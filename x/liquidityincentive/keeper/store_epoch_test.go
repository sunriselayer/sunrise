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
		id, _ := keeper.AppendEpoch(ctx, items[i])
		items[i].Id = id
	}
	return items
}

func TestEpochGet(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	for _, item := range items {
		got, found, err := f.keeper.GetEpoch(f.ctx, item.Id)
		require.NoError(t, err)
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
		err := f.keeper.RemoveEpoch(f.ctx, item.Id)
		require.NoError(t, err)
		_, found, err := f.keeper.GetEpoch(f.ctx, item.Id)
		require.NoError(t, err)
		require.False(t, found)
	}
}

func TestEpochGetAll(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	epochs, err := f.keeper.GetAllEpoch(f.ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(epochs),
	)
}

func TestEpochCount(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	count, err := f.keeper.GetEpochCount(f.ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(len(items)), count)
}

func TestGetLastEpoch(t *testing.T) {
	f := initFixture(t)
	items := createNEpoch(f.keeper, f.ctx, 10)
	got, found, err := f.keeper.GetLastEpoch(f.ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&items[len(items)-1]),
		nullify.Fill(&got),
	)
}
