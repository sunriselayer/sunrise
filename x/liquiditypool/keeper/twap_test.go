package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/sunriselayer/sunrise-app/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
	keepertest "github.com/sunriselayer/sunrise-app/testutil/keeper"
	"github.com/sunriselayer/sunrise-app/testutil/nullify"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNTwap(keeper keeper.Keeper, ctx context.Context, n int) []types.Twap {
	items := make([]types.Twap, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
        
		keeper.SetTwap(ctx, items[i])
	}
	return items
}

func TestTwapGet(t *testing.T) {
	keeper, ctx := keepertest.LiquiditypoolKeeper(t)
	items := createNTwap(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTwap(ctx,
		    item.Index,
            
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
		    item.Index,
            
		)
		_, found := keeper.GetTwap(ctx,
		    item.Index,
            
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
