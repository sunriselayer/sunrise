package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
	"github.com/sunriselayer/sunrise-app/testutil/nullify"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPair(keeper keeper.Keeper, ctx context.Context, n int) []types.Pair {
	items := make([]types.Pair, n)
	for i := range items {
		index := strconv.Itoa(i)
		items[i].BaseDenom = "base" + index
		items[i].QuoteDenom = "quote" + index

		keeper.SetPair(ctx, items[i])
	}
	return items
}

func TestPairGet(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPair(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPair(ctx,
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
func TestPairRemove(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPair(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePair(ctx,
			item.BaseDenom,
			item.QuoteDenom,
		)
		_, found := keeper.GetPair(ctx,
			item.BaseDenom,
			item.QuoteDenom,
		)
		require.False(t, found)
	}
}

func TestPairGetAll(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNPair(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPair(ctx)),
	)
}
