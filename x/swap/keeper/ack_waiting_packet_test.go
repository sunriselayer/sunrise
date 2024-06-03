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

func createNAckWaitingPacket(keeper keeper.Keeper, ctx context.Context, n int) []types.AckWaitingPacket {
	items := make([]types.AckWaitingPacket, n)
	for i := range items {
		items[i].Index = types.PacketIndex{
			PortId:    strconv.Itoa(i),
			ChannelId: strconv.Itoa(i),
			Sequence:  uint64(i),
		}

		keeper.SetAckWaitingPacket(ctx, items[i])
	}
	return items
}

func TestAckWaitingPacketGet(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNAckWaitingPacket(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetAckWaitingPacket(ctx,
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
func TestAckWaitingPacketRemove(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNAckWaitingPacket(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAckWaitingPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		_, found := keeper.GetAckWaitingPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		require.False(t, found)
	}
}

func TestAckWaitingPacketGetAll(t *testing.T) {
	keeper, ctx := keepertest.SwapKeeper(t)
	items := createNAckWaitingPacket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAckWaitingPacket(ctx)),
	)
}
