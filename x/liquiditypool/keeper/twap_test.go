package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNTwap(keeper keeper.Keeper, ctx context.Context, n int) []types.Twap {
	items := make([]types.Twap, n)
	for i := range items {
		index := strconv.Itoa(i)
		items[i].BaseDenom = "base" + index
		items[i].QuoteDenom = "quote" + index

		keeper.SetTwap(ctx, items[i])
	}
	return items
}

func TestTwapGet(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNTwap(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTwap(ctx,
			item.BaseDenom,
			item.QuoteDenom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestTwapRemove(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNTwap(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTwap(ctx,
			item.BaseDenom,
			item.QuoteDenom,
		)
		_, found := keeper.GetTwap(ctx,
			item.BaseDenom,
			item.QuoteDenom,
		)
		require.False(t, found)
	}
}

func TestTwapGetAll(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNTwap(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllTwap(ctx)),
	)
}
