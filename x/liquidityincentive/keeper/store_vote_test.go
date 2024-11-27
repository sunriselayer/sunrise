package keeper_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNVote(keeper keeper.Keeper, ctx context.Context, n int) []types.Vote {
	items := make([]types.Vote, n)
	for i := range items {
		items[i].Sender = sdk.AccAddress(fmt.Sprintf("sender%d", i)).String()
		items[i].Weights = []types.PoolWeight{
			{
				PoolId: 1,
				Weight: math.LegacyOneDec(),
			},
		}

		keeper.SetVote(ctx, items[i])
	}
	return items
}

func TestVoteSet(t *testing.T) {
	keeper, _, ctx := keepertest.LiquidityincentiveKeeper(t)
	keeper.SetVote(ctx, types.Vote{
		Sender:  "sender1",
		Weights: []types.PoolWeight{{PoolId: 1, Weight: math.LegacyOneDec()}, {PoolId: 2, Weight: math.LegacyOneDec()}},
	})
	keeper.SetVote(ctx, types.Vote{
		Sender:  "sender2",
		Weights: []types.PoolWeight{{PoolId: 1, Weight: math.LegacyOneDec()}},
	})
	require.ElementsMatch(t,
		nullify.Fill([]types.Vote{{
			Sender:  "sender1",
			Weights: []types.PoolWeight{{PoolId: 1, Weight: math.LegacyOneDec()}, {PoolId: 2, Weight: math.LegacyOneDec()}},
		}, {
			Sender:  "sender2",
			Weights: []types.PoolWeight{{PoolId: 1, Weight: math.LegacyOneDec()}},
		}}),
		nullify.Fill(keeper.GetAllVotes(ctx)),
	)
}

func TestVoteGet(t *testing.T) {
	keeper, _, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNVote(keeper, ctx, 10)
	for i, item := range items {
		address := sdk.AccAddress(fmt.Sprintf("sender%d", i)).String()
		rst, found := keeper.GetVote(ctx, address)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestVoteRemove(t *testing.T) {
	keeper, _, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNVote(keeper, ctx, 10)
	for i := range items {
		address := sdk.AccAddress(fmt.Sprintf("sender%d", i)).String()
		keeper.RemoveVote(ctx, address)
		_, found := keeper.GetVote(ctx, address)
		require.False(t, found)
	}
}

func TestVoteGetAll(t *testing.T) {
	keeper, _, ctx := keepertest.LiquidityincentiveKeeper(t)
	items := createNVote(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllVotes(ctx)),
	)
}
