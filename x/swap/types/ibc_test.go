package types

import (
	"strings"
	"testing"
	time "time"

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
			AmountStrategy: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				},
			},
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Next:     "{\"wasm\":{\"contract\":\"neutron13lw4uyaxc09d3qvgunc8crtxcvwd8pn5xzx2xzlw53nsfgq3y8ps4q2dr5\",\"msg\":{\"send_to_evm\":{\"destination_chain\":\"ethereum-sepolia\",\"destination_contract\":\"0x8ef2c2b9825a52c44bff05b4dd7b72899ccbd4e4\",\"recipient\":\"0x4793755541ae9f950a68fc7fc2b3bd2cc9397b9a\"}}}}",
			},
		},
	}

	m := jsonpb.Marshaler{OrigName: true}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2,"next":"{\"wasm\":{\"contract\":\"neutron13lw4uyaxc09d3qvgunc8crtxcvwd8pn5xzx2xzlw53nsfgq3y8ps4q2dr5\",\"msg\":{\"send_to_evm\":{\"destination_chain\":\"ethereum-sepolia\",\"destination_contract\":\"0x8ef2c2b9825a52c44bff05b4dd7b72899ccbd4e4\",\"recipient\":\"0x4793755541ae9f950a68fc7fc2b3bd2cc9397b9a\"}}}}"}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
	require.Equal(t, metadata, &packetMetadata)
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
			AmountStrategy: &SwapMetadata_ExactAmountIn{
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
			AmountStrategy: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				}},
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
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
			AmountStrategy: &SwapMetadata_ExactAmountOut{
				ExactAmountOut: &ExactAmountOut{
					AmountOut: sdkmath.NewInt(1000),
					Change: &ForwardMetadata{
						Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
						Port:     "transfer",
						Channel:  "channel-2",
						Retries:  retries,
					},
				}},
			Forward: &ForwardMetadata{
				Receiver: "cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9",
				Port:     "transfer",
				Channel:  "channel-2",
				Retries:  retries,
				Timeout:  time.Duration(time.Hour),
			},
		},
	}

	m := jsonpb.Marshaler{OrigName: true}
	js, err := m.MarshalToString(&packetMetadata)
	require.NoError(t, err)

	require.Equal(t, js, `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_out":{"amount_out":"1000","change":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","retries":2}},"forward":{"receiver":"cosmos1qnk2n4nlkpw9xfqntladh74w6ujtulwn7j8za9","port":"transfer","channel":"channel-2","timeout":"3600s","retries":2}}}`)

	metadata := &PacketMetadata{}
	err = jsonpb.Unmarshal(strings.NewReader(js), metadata)
	require.NoError(t, err)
}

func TestDecodePacketMetadata_NoForward(t *testing.T) {
	memoString := `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_in":{"min_amount_out":"1"}}}`
	m, err := DecodeSwapMetadata(memoString)
	require.NoError(t, err)

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
			AmountStrategy: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				},
			},
		},
	}

	require.Equal(t, *m, packetMetadata)
}

func TestDecodePacketMetadata_NoNext(t *testing.T) {
	memoString := `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"neutron1s2gtqhnj9d6q5wjr44ll6uyd3xwn9a7fcn8t53yewjtq04ru52fsgupa3j","port":"transfer","channel":"channel-1","timeout":"3600s","retries":2}}}`
	m, err := DecodeSwapMetadata(memoString)
	require.NoError(t, err)

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
			AmountStrategy: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				},
			},
			Forward: &ForwardMetadata{
				Receiver: "neutron1s2gtqhnj9d6q5wjr44ll6uyd3xwn9a7fcn8t53yewjtq04ru52fsgupa3j",
				Port:     "transfer",
				Channel:  "channel-1",
				Timeout:  3600000000000,
				Retries:  retries,
			},
		},
	}

	require.Equal(t, *m, packetMetadata)
}

func TestDecodePacketMetadata_ForwardAndNext(t *testing.T) {
	memoString := `{"swap":{"interface_provider":"sunrise18atdu5vvsg95sdpvdwsv7kevlzg8jhtuk7hs4y","route":{"denom_in":"tokenIn","denom_out":"tokenOut","pool":{"pool_id":"1"}},"exact_amount_in":{"min_amount_out":"1"},"forward":{"receiver":"neutron1s2gtqhnj9d6q5wjr44ll6uyd3xwn9a7fcn8t53yewjtq04ru52fsgupa3j","port":"transfer","channel":"channel-1","timeout":"3600s","retries":2,"next":{"wasm":{"contract":"neutron1s2gtqhnj9d6q5wjr44ll6uyd3xwn9a7fcn8t53yewjtq04ru52fsgupa3j","msg":{"send_to_evm":{"destination_chain":"ethereum-sepolia","destination_contract":"0x8ef2c2b9825a52c44bff05b4dd7b72899ccbd4e4","recipient":"0x4793755541ae9f950a68fc7fc2b3bd2cc9397b9a","fee":"277912"}}}}}}}`
	m, err := DecodeSwapMetadata(memoString)
	require.NoError(t, err)

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
			AmountStrategy: &SwapMetadata_ExactAmountIn{
				ExactAmountIn: &ExactAmountIn{
					MinAmountOut: sdkmath.OneInt(),
				},
			},
			Forward: &ForwardMetadata{
				Receiver: "neutron1s2gtqhnj9d6q5wjr44ll6uyd3xwn9a7fcn8t53yewjtq04ru52fsgupa3j",
				Port:     "transfer",
				Channel:  "channel-1",
				Timeout:  3600000000000,
				Retries:  retries,
				Next:     "{\"wasm\":{\"contract\":\"neutron1s2gtqhnj9d6q5wjr44ll6uyd3xwn9a7fcn8t53yewjtq04ru52fsgupa3j\",\"msg\":{\"send_to_evm\":{\"destination_chain\":\"ethereum-sepolia\",\"destination_contract\":\"0x8ef2c2b9825a52c44bff05b4dd7b72899ccbd4e4\",\"fee\":\"277912\",\"recipient\":\"0x4793755541ae9f950a68fc7fc2b3bd2cc9397b9a\"}}}}",
			},
		},
	}

	require.Equal(t, *m, packetMetadata)
}
