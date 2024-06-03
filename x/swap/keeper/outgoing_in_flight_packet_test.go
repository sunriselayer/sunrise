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

func createNOutgoingInFlightPacket(keeper keeper.Keeper, ctx context.Context, n int) []types.OutgoingInFlightPacket {
	items := make([]types.OutgoingInFlightPacket, n)
	for i := range items {
		items[i].Index.PortId = strconv.Itoa(i)
		items[i].Index.ChannelId = strconv.Itoa(i)
		items[i].Index.Sequence = uint64(i)

		keeper.SetOutgoingInFlightPacket(ctx, items[i])
	}
	return items
}

func TestOutgoingInFlightPacketGet(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNOutgoingInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetOutgoingInFlightPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestOutgoingInFlightPacketRemove(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNOutgoingInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveOutgoingInFlightPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		_, found := keeper.GetOutgoingInFlightPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		require.False(t, found)
	}
}

func TestOutgoingInFlightPacketGetAll(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNOutgoingInFlightPacket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOutgoingInFlightPacket(ctx)),
	)
}
