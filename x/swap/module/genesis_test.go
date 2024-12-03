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

		IncomingInFlightPackets: []types.IncomingInFlightPacket{
			{
				Index: types.PacketIndex{
					PortId:    "0",
					ChannelId: "0",
					Sequence:  0,
				},
			},
			{
				Index: types.PacketIndex{
					PortId:    "1",
					ChannelId: "1",
					Sequence:  1,
				},
			},
		},
		OutgoingInFlightPackets: []types.OutgoingInFlightPacket{
			{
				Index: types.PacketIndex{
					PortId:    "0",
					ChannelId: "0",
					Sequence:  0,
				},
			},
			{
				Index: types.PacketIndex{
					PortId:    "1",
					ChannelId: "1",
					Sequence:  1,
				},
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, _, ctx := keepertest.SwapKeeper(t)
	swap.InitGenesis(ctx, k, genesisState)
	got := swap.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.OutgoingInFlightPackets, got.OutgoingInFlightPackets)
	require.ElementsMatch(t, genesisState.IncomingInFlightPackets, got.IncomingInFlightPackets)
	// this line is used by starport scaffolding # genesis/test/assert
}
