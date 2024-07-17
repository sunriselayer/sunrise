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
			ExactAmountIn: &ExactAmountIn{
				MinAmountOut: sdkmath.OneInt(),
			},
			ExactAmountOut: nil,
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Next:     nil,
			},
		},
	}

	m := jsonpb.Marshaler{}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interfaceProvider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denomIn":"tokenIn","denomOut":"tokenOut","pool":{"poolId":"1"}},"exactAmountIn":{"minAmountOut":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

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
			ExactAmountIn: &ExactAmountIn{
				MinAmountOut: sdkmath.OneInt(),
			},
			ExactAmountOut: nil,
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Next:     nil,
			},
		},
	}

	m := jsonpb.Marshaler{}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interfaceProvider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denomIn":"tokenIn","denomOut":"tokenOut","series":{"routes":[{"denomIn":"tokenIn","denomOut":"tokenOut","pool":{"poolId":"1"}},{"denomIn":"tokenOut","denomOut":"tokenOut2","pool":{"poolId":"2"}}]}},"exactAmountIn":{"minAmountOut":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

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
			ExactAmountIn: nil,
			ExactAmountOut: &ExactAmountOut{
				AmountOut: sdkmath.NewInt(1000),
				Change: &ForwardMetadata{
					Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
					Port:     "transfer",
					Channel:  "channel-2",
					Retries:  retries,
					Next:     nil,
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

	m := jsonpb.Marshaler{}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interfaceProvider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denomIn":"tokenIn","denomOut":"tokenOut","pool":{"poolId":"1"}},"exactAmountOut":{"amountOut":"1000","change":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
}
