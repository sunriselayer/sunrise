package types

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	"github.com/stretchr/testify/require"
)

func TestEncodePacketMetadata_ExactAmountIn(t *testing.T) {
	retries := uint8(2)
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: SwapRoute{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut",
				Strategy: SwapStrategy{
					Pool: &RoutePool{
						PoolId: 1,
					},
				},
			},
			ExactAmountIn: &ExactAmountIn{
				MinAmountOut: sdkmath.OneInt(),
			},
			ExactAmountOut: nil,
			Forward: &packetforwardtypes.ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  &retries,
				Next:     nil,
			},
		},
	}

	metadataJson, err := json.Marshal(packetMetadata)
	require.NoError(t, err)

	require.Equal(t, string(metadataJson), `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","strategy":{"pool":{"pool_id":1}}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := PacketMetadata{}
	err = json.Unmarshal(metadataJson, &metadata)
	require.NoError(t, err)
}

func TestEncodePacketMetadata_ExactAmountInSeries(t *testing.T) {
	retries := uint8(2)
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: SwapRoute{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut2",
				Strategy: SwapStrategy{
					Series: &RouteSeries{
						Routes: []Route{{
							DenomIn:  "tokenIn",
							DenomOut: "tokenOut",
							Strategy: &Route_Pool{
								Pool: &RoutePool{
									PoolId: 1,
								},
							}},
							{
								DenomIn:  "tokenOut",
								DenomOut: "tokenOut2",
								Strategy: &Route_Pool{
									Pool: &RoutePool{
										PoolId: 2,
									},
								}},
						},
					}},
			},
			ExactAmountIn: &ExactAmountIn{
				MinAmountOut: sdkmath.OneInt(),
			},
			ExactAmountOut: nil,
			Forward: &packetforwardtypes.ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  &retries,
				Next:     nil,
			},
		},
	}

	metadataJson, err := json.Marshal(packetMetadata)
	require.NoError(t, err)

	require.Equal(t, string(metadataJson), `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut2","strategy":{"series":{"routes":[{"denom_in":"tokenIn","denom_out":"tokenOut","Strategy":{"pool":{"pool_id":1}}},{"denom_in":"tokenOut","denom_out":"tokenOut2","Strategy":{"pool":{"pool_id":2}}}]}}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := PacketMetadata{}
	err = json.Unmarshal(metadataJson, &metadata)
	require.NoError(t, err)
}

func TestEncodePacketMetadata_ExactAmountOut(t *testing.T) {
	retries := uint8(2)
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: SwapRoute{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut",
				Strategy: SwapStrategy{
					Pool: &RoutePool{
						PoolId: 1,
					},
				},
			},
			ExactAmountIn: nil,
			ExactAmountOut: &ExactAmountOut{
				AmountOut: sdkmath.NewInt(1000),
				Change: &packetforwardtypes.ForwardMetadata{
					Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
					Port:     "transfer",
					Channel:  "channel-2",
					Retries:  &retries,
					Next:     nil,
				},
			},
			Forward: &packetforwardtypes.ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  &retries,
				Next:     nil,
			},
		},
	}

	metadataJson, err := json.Marshal(packetMetadata)
	require.NoError(t, err)

	require.Equal(t, string(metadataJson), `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","strategy":{"pool":{"pool_id":1}}},"exact_amount_out":{"amount_out":"1000","change":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := PacketMetadata{}
	err = json.Unmarshal(metadataJson, &metadata)
	require.NoError(t, err)
}
