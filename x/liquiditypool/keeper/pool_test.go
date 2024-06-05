package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func createNPool(keeper keeper.Keeper, ctx context.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		items[i].Id = keeper.AppendPool(ctx, items[i])
	}
	return items
}

func TestPoolGet(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPool(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePool(ctx, item.Id)
		_, found := keeper.GetPool(ctx, item.Id)
		require.False(t, found)
	}
}

func TestPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPools(ctx)),
	)
}

func TestPoolCount(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPool(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPoolCount(ctx))
}
