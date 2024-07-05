package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func TestSwapExactAmountOut(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"
	senderAcc := sdk.MustAccAddressFromBech32(sender)
	tests := []struct {
		desc              string
		interfaceProvider string
		route             types.Route
		maxAmountIn       math.Int
		amountOut         math.Int
		expResult         types.RouteResult
		expInterfaceFee   math.Int
		expErr            error
	}{
		{
			desc:              "Single pool route testing",
			interfaceProvider: "",
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Pool{
					Pool: &types.RoutePool{
						PoolId: 0,
					},
				},
			},
			maxAmountIn: math.OneInt(),
			amountOut:   math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Pool{
					Pool: &types.RouteResultPool{
						PoolId: 0,
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:              "Route Series testing",
			interfaceProvider: "",
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Series{
					Series: &types.RouteSeries{
						Routes: []types.Route{
							{
								DenomIn:  "base",
								DenomOut: "quote",
								Strategy: &types.Route_Pool{
									Pool: &types.RoutePool{
										PoolId: 0,
									},
								},
							},
						},
					},
				},
			},
			maxAmountIn: math.OneInt(),
			amountOut:   math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Series{
					Series: &types.RouteResultSeries{
						RouteResults: []types.RouteResult{
							{
								TokenIn:  sdk.NewInt64Coin("base", 0),
								TokenOut: sdk.NewInt64Coin("quote", 1),
								Strategy: &types.RouteResult_Pool{
									Pool: &types.RouteResultPool{
										PoolId: 0,
									},
								},
							},
						},
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:              "Route Parallel testing",
			interfaceProvider: "",
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Parallel{
					Parallel: &types.RouteParallel{
						Routes: []types.Route{
							{
								DenomIn:  "base",
								DenomOut: "quote",
								Strategy: &types.Route_Pool{
									Pool: &types.RoutePool{
										PoolId: 0,
									},
								},
							},
						},
						Weights: []math.LegacyDec{math.LegacyOneDec()},
					},
				},
			},
			maxAmountIn: math.OneInt(),
			amountOut:   math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Parallel{
					Parallel: &types.RouteResultParallel{
						RouteResults: []types.RouteResult{
							{
								TokenIn:  sdk.NewInt64Coin("base", 0),
								TokenOut: sdk.NewInt64Coin("quote", 0),
								Strategy: &types.RouteResult_Pool{
									Pool: &types.RouteResultPool{
										PoolId: 0,
									},
								},
							},
						},
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:              "Empty route object route",
			interfaceProvider: "",
			route:             types.Route{},
			maxAmountIn:       math.OneInt(),
			amountOut:         math.OneInt(),
			expResult:         types.RouteResult{},
			expInterfaceFee:   math.ZeroInt(),
			expErr:            types.UnknownStrategyType,
		},
		{
			desc:              "nil strategy route",
			interfaceProvider: "",
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: nil,
			},
			maxAmountIn:     math.OneInt(),
			amountOut:       math.OneInt(),
			expResult:       types.RouteResult{},
			expInterfaceFee: math.ZeroInt(),
			expErr:          types.UnknownStrategyType,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			keeper, mocks, ctx := keepertest.SwapKeeper(t)
			mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), gomock.Any()).Return(liquiditypooltypes.Pool{}, true).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().SwapExactAmountOut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(math.OneInt(), nil).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().CalculateResultExactAmountOut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(math.OneInt(), nil).AnyTimes()

			result, interfaceFee, err := keeper.SwapExactAmountOut(ctx, senderAcc, tc.interfaceProvider, tc.route, tc.maxAmountIn, tc.amountOut)
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, result.TokenIn.String(), tc.expResult.TokenIn.String())
				require.Equal(t, result.TokenOut.String(), tc.expResult.TokenOut.String())
				require.Equal(t, interfaceFee.String(), tc.expInterfaceFee.String())
			}
		})
	}
}

func TestCalculateResultExactAmountOut(t *testing.T) {
	tests := []struct {
		desc            string
		interfaceFee    bool
		route           types.Route
		amountOut       math.Int
		expResult       types.RouteResult
		expInterfaceFee math.Int
		expErr          error
	}{
		{
			desc:         "Single pool route testing",
			interfaceFee: false,
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Pool{
					Pool: &types.RoutePool{
						PoolId: 0,
					},
				},
			},
			amountOut: math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Pool{
					Pool: &types.RouteResultPool{
						PoolId: 0,
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:         "Positive interface fee testing",
			interfaceFee: true,
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Pool{
					Pool: &types.RoutePool{
						PoolId: 0,
					},
				},
			},
			amountOut: math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Pool{
					Pool: &types.RouteResultPool{
						PoolId: 0,
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:         "Route Series testing",
			interfaceFee: false,
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Series{
					Series: &types.RouteSeries{
						Routes: []types.Route{
							{
								DenomIn:  "base",
								DenomOut: "quote",
								Strategy: &types.Route_Pool{
									Pool: &types.RoutePool{
										PoolId: 0,
									},
								},
							},
						},
					},
				},
			},
			amountOut: math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Series{
					Series: &types.RouteResultSeries{
						RouteResults: []types.RouteResult{
							{
								TokenIn:  sdk.NewInt64Coin("base", 0),
								TokenOut: sdk.NewInt64Coin("quote", 1),
								Strategy: &types.RouteResult_Pool{
									Pool: &types.RouteResultPool{
										PoolId: 0,
									},
								},
							},
						},
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:         "Route Parallel testing",
			interfaceFee: false,
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: &types.Route_Parallel{
					Parallel: &types.RouteParallel{
						Routes: []types.Route{
							{
								DenomIn:  "base",
								DenomOut: "quote",
								Strategy: &types.Route_Pool{
									Pool: &types.RoutePool{
										PoolId: 0,
									},
								},
							},
						},
						Weights: []math.LegacyDec{math.LegacyOneDec()},
					},
				},
			},
			amountOut: math.OneInt(),
			expResult: types.RouteResult{
				TokenIn:  sdk.NewInt64Coin("base", 1),
				TokenOut: sdk.NewInt64Coin("quote", 1),
				Strategy: &types.RouteResult_Parallel{
					Parallel: &types.RouteResultParallel{
						RouteResults: []types.RouteResult{
							{
								TokenIn:  sdk.NewInt64Coin("base", 0),
								TokenOut: sdk.NewInt64Coin("quote", 0),
								Strategy: &types.RouteResult_Pool{
									Pool: &types.RouteResultPool{
										PoolId: 0,
									},
								},
							},
						},
					},
				},
			},
			expInterfaceFee: math.ZeroInt(),
			expErr:          nil,
		},
		{
			desc:            "Empty route object route",
			interfaceFee:    false,
			route:           types.Route{},
			amountOut:       math.OneInt(),
			expResult:       types.RouteResult{},
			expInterfaceFee: math.ZeroInt(),
			expErr:          types.UnknownStrategyType,
		},
		{
			desc:         "nil strategy route",
			interfaceFee: false,
			route: types.Route{
				DenomIn:  "base",
				DenomOut: "quote",
				Strategy: nil,
			},
			amountOut:       math.OneInt(),
			expResult:       types.RouteResult{},
			expInterfaceFee: math.ZeroInt(),
			expErr:          types.UnknownStrategyType,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			keeper, mocks, ctx := keepertest.SwapKeeper(t)
			mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), gomock.Any()).Return(liquiditypooltypes.Pool{}, true).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().SwapExactAmountOut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(math.OneInt(), nil).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().CalculateResultExactAmountOut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(math.NewInt(1), nil).AnyTimes()

			result, interfaceFee, err := keeper.CalculateResultExactAmountOut(ctx, tc.interfaceFee, tc.route, tc.amountOut)
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)

				require.Equal(t, result.TokenIn.String(), tc.expResult.TokenIn.String())
				require.Equal(t, interfaceFee.String(), tc.expInterfaceFee.String())
			}
		})
	}
}

// TODO: add test for generateResultExactAmountOut
// TODO: add test for calculateResultRoutePoolExactAmountOut
// TODO: add test for swapRouteExactAmountOut
// TODO: add test for swapRoutePoolExactAmountOut
