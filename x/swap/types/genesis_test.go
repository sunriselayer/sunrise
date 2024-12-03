package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated incomingInFlightPacket",
			genState: &types.GenesisState{
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
							PortId:    "0",
							ChannelId: "0",
							Sequence:  0,
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated outgoingInFlightPacket",
			genState: &types.GenesisState{
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
							PortId:    "0",
							ChannelId: "0",
							Sequence:  0,
						},
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
