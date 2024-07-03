package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
	"go.uber.org/mock/gomock"
)

func TestCalculationSwapExactAmountIn(t *testing.T) {
	tests := []struct {
		desc            string
		interfaceFee    bool
		route           types.Route
		amountIn        math.Int
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
			amountIn: math.OneInt(),
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
			amountIn: math.OneInt(),
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
			amountIn: math.OneInt(),
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
			amountIn: math.OneInt(),
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
			amountIn:        math.ZeroInt(),
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
			amountIn:        math.ZeroInt(),
			expResult:       types.RouteResult{},
			expInterfaceFee: math.ZeroInt(),
			expErr:          types.UnknownStrategyType,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			keeper, mocks, ctx := keepertest.SwapKeeper(t)
			mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), gomock.Any()).Return(liquiditypooltypes.Pool{}, true).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().CalculateResultExactAmountIn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(math.OneInt(), nil).AnyTimes()

			resp, err := keeper.CalculationSwapExactAmountIn(ctx, &types.QueryCalculationSwapExactAmountInRequest{
				HasInterfaceFee: tc.interfaceFee,
				Route:           &tc.route,
				AmountIn:        tc.amountIn.String(),
			})
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)

				require.Equal(t, resp.Result.TokenOut.String(), tc.expResult.TokenOut.String())
				require.Equal(t, resp.InterfaceProviderFee.String(), tc.expInterfaceFee.String())
				require.Equal(t, resp.AmountOut.String(), tc.expResult.TokenOut.Amount.Sub(tc.expInterfaceFee).String())
			}
		})
	}
}

func TestCalculationSwapExactAmountOut(t *testing.T) {
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

			resp, err := keeper.CalculationSwapExactAmountOut(ctx, &types.QueryCalculationSwapExactAmountOutRequest{
				HasInterfaceFee: tc.interfaceFee,
				Route:           &tc.route,
				AmountOut:       tc.amountOut.String(),
			})
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)

				require.Equal(t, resp.Result.TokenIn.String(), tc.expResult.TokenIn.String())
				require.Equal(t, resp.InterfaceProviderFee.String(), tc.expInterfaceFee.String())
			}
		})
	}
}
