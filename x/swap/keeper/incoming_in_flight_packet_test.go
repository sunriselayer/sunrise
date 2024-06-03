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

func createNIncomingInFlightPacket(keeper keeper.Keeper, ctx context.Context, n int) []types.IncomingInFlightPacket {
	items := make([]types.IncomingInFlightPacket, n)
	for i := range items {
		items[i].Index = types.PacketIndex{
			PortId:    strconv.Itoa(i),
			ChannelId: strconv.Itoa(i),
			Sequence:  uint64(i),
		}

		keeper.SetIncomingInFlightPacket(ctx, items[i])
	}
	return items
}

func TestIncomingInFlightPacketGet(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNIncomingInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetIncomingInFlightPacket(ctx,
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
func TestIncomingInFlightPacketRemove(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNIncomingInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveIncomingInFlightPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		_, found := keeper.GetIncomingInFlightPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		require.False(t, found)
	}
}

func TestIncomingInFlightPacketGetAll(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNIncomingInFlightPacket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllIncomingInFlightPacket(ctx)),
	)
}
