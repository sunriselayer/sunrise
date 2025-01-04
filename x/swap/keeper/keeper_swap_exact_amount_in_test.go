package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

func TestSwapExactAmountIn(t *testing.T) {
	sender := "sunrise126ss57ayztn5287spvxq0dpdfarj6rk0v3p06f"
	senderAcc := sdk.MustAccAddressFromBech32(sender)
	tests := []struct {
		desc              string
		interfaceProvider string
		route             types.Route
		amountIn          math.Int
		minAmountOut      math.Int
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
			amountIn:     math.OneInt(),
			minAmountOut: math.ZeroInt(),
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
			amountIn:     math.OneInt(),
			minAmountOut: math.ZeroInt(),
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
			amountIn:     math.OneInt(),
			minAmountOut: math.ZeroInt(),
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
			amountIn:          math.ZeroInt(),
			minAmountOut:      math.ZeroInt(),
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
			amountIn:        math.ZeroInt(),
			minAmountOut:    math.ZeroInt(),
			expResult:       types.RouteResult{},
			expInterfaceFee: math.ZeroInt(),
			expErr:          types.UnknownStrategyType,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			f := initFixture(t)
			ctx := sdk.UnwrapSDKContext(f.ctx)
			keeper := f.keeper
			mocks := getMocks(t)

			mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), gomock.Any()).Return(liquiditypooltypes.Pool{}, true).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().SwapExactAmountIn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(math.OneInt(), nil).AnyTimes()

			result, interfaceFee, err := keeper.SwapExactAmountIn(ctx, senderAcc, tc.interfaceProvider, tc.route, tc.amountIn, tc.minAmountOut)
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

// TODO: add test for calculateInterfaceFeeExactAmountIn
// TODO: add test for CalculateResultExactAmountIn
// TODO: add test for generateResultExactAmountIn
// TODO: add test for calculateResultRouteExactAmountIn
// TODO: add test for calculateResultRoutePoolExactAmountIn
// TODO: add test for swapRouteExactAmountIn
// TODO: add test for swapRoutePoolExactAmountIn
