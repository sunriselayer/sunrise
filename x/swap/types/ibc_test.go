package types

import (
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/require"
)

func TestEncodePacketMetadata_ExactAmountIn(t *testing.T) {
	retries := uint32(2)
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: &Route{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut",
				Strategy: &Route_Pool{
					Pool: &RoutePool{
						PoolId: 1,
					},
				},
			},
			SwapType: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				},
			},
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Next:     nil,
			},
		},
	}

	m := jsonpb.Marshaler{OrigName: true}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
}

func TestEncodePacketMetadataNoForward_ExactAmountIn(t *testing.T) {
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: &Route{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut",
				Strategy: &Route_Pool{
					Pool: &RoutePool{
						PoolId: 1,
					},
				},
			},
			SwapType: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				},
			},
		},
	}

	m := jsonpb.Marshaler{OrigName: true}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_in":{"min_amount_out":"1"}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
}

func TestEncodePacketMetadata_ExactAmountInSeries(t *testing.T) {
	retries := uint32(2)
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: &Route{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut",
				Strategy: &Route_Series{
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
					},
				},
			},
			SwapType: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				}},
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Next:     nil,
			},
		},
	}

	m := jsonpb.Marshaler{OrigName: true}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","series":{"routes":[{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},{"denom_in":"tokenOut","denom_out":"tokenOut2","pool":{"pool_id":"2"}}]}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
}

func TestEncodePacketMetadata_ExactAmountOut(t *testing.T) {
	retries := uint32(2)
	packetMetadata := PacketMetadata{
		Swap: &SwapMetadata{
			InterfaceProvider: "sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y",
			Route: &Route{
				DenomIn:  "tokenIn",
				DenomOut: "tokenOut",
				Strategy: &Route_Pool{
					Pool: &RoutePool{
						PoolId: 1,
					},
				},
			},
			SwapType: &SwapMetadata_ExactAmountOut{
				ExactAmountOut: &ExactAmountOut{
					AmountOut: sdkmath.NewInt(1000),
					Change: &ForwardMetadata{
						Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
						Port:     "transfer",
						Channel:  "channel-2",
						Retries:  retries,
						Next:     nil,
					},
				}},
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Next:     nil,
			},
		},
	}

	m := jsonpb.Marshaler{OrigName: true}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_out":{"amount_out":"1000","change":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
}
