package swap_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/testutil/nullify"
	swap "github.com/sunriselayer/sunrise/x/swap/module"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		InFlightPacketList: []types.InFlightPacket{
			{
				Index: types.PacketIndex{
					SrcPortId:    "0",
					SrcChannelId: "0",
					Sequence:     0,
				},
			},
			{
				Index: types.PacketIndex{
					SrcPortId:    "1",
					SrcChannelId: "1",
					Sequence:     1,
				},
			},
		},
		AckWaitingPacketList: []types.AckWaitingPacket{
			{
				Index: types.PacketIndex{
					SrcPortId:    "0",
					SrcChannelId: "0",
					Sequence:     0,
				},
			},
			{
				Index: types.PacketIndex{
					SrcPortId:    "1",
					SrcChannelId: "1",
					Sequence:     1,
				},
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SwapKeeper(t)
	swap.InitGenesis(ctx, k, genesisState)
	got := swap.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.InFlightPacketList, got.InFlightPacketList)
	require.ElementsMatch(t, genesisState.AckWaitingPacketList, got.AckWaitingPacketList)
	// this line is used by starport scaffolding # genesis/test/assert
}
