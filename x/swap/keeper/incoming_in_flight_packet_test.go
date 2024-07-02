package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/math"
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
		items[i] = types.IncomingInFlightPacket{
			Index: types.PacketIndex{
				PortId:    strconv.Itoa(i),
				ChannelId: strconv.Itoa(i),
				Sequence:  uint64(i),
			},
			Data:             []byte{},
			SrcPortId:        "transfer",
			SrcChannelId:     "channel-1",
			TimeoutHeight:    "",
			TimeoutTimestamp: 1717912068,
			Ack:              []byte{},
			Result:           types.RouteResult{},
			InterfaceFee:     math.NewInt(1),
			Change:           nil,
			Forward:          nil,
		}
		keeper.SetIncomingInFlightPacket(ctx, items[i])
	}
	return items
}

func TestIncomingInFlightPacketGet(t *testing.T) {
	keeper, _, ctx := keepertest.SwapKeeper(t)
	items := createNIncomingInFlightPacket(keeper, ctx, 10)
	for _, item := range items {
		_, found := keeper.GetIncomingInFlightPacket(ctx,
			item.Index.PortId,
			item.Index.ChannelId,
			item.Index.Sequence,
		)
		require.True(t, found)
	}
}
func TestIncomingInFlightPacketRemove(t *testing.T) {
	keeper, _, ctx := keepertest.SwapKeeper(t)
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
	keeper, _, ctx := keepertest.SwapKeeper(t)
	items := createNIncomingInFlightPacket(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetIncomingInFlightPackets(ctx)),
	)
}
