package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	"github.com/sunriselayer/sunrise/x/swap/keeper"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNInFlightPacket(keeper keeper.Keeper, ctx context.Context, n int) []types.InFlightPacket {
	items := make([]types.InFlightPacket, n)
	for i := range items {
		items[i].Index.SrcPortId = strconv.Itoa(i)
		items[i].Index.SrcChannelId = strconv.Itoa(i)
		items[i].Index.Sequence = uint64(i)

		keeper.SetInFlightPacket(ctx, items[i])
	}
	return items
}

func TestInFlightPacketGet(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetInFlightPacket(ctx,
			item.Index.SrcPortId,
			item.Index.SrcChannelId,
			item.Index.Sequence,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestInFlightPacketRemove(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveInFlightPacket(ctx,
			item.Index.SrcPortId,
			item.Index.SrcChannelId,
			item.Index.Sequence,
		)
		_, found := keeper.GetInFlightPacket(ctx,
			item.Index.SrcPortId,
			item.Index.SrcChannelId,
			item.Index.Sequence,
		)
		require.False(t, found)
	}
}

func TestInFlightPacketGetAll(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNInFlightPacket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllInFlightPacket(ctx)),
	)
}
