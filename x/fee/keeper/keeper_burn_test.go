package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/fee/types"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
)

const (
	feeDenom  = "fee"
	burnDenom = "burn"
)

func TestBurn(t *testing.T) {
	fee := sdk.NewInt64Coin(feeDenom, 1000)
	fees := sdk.NewCoins(fee)

	testCases := []struct {
		name        string
		setup       func(f *fixture)
		fees        sdk.Coins
		params      types.Params
		expectError bool
	}{
		{
			name: "successful burn without swap",
			fees: fees,
			params: types.Params{
				FeeDenom:  feeDenom,
				BurnDenom: feeDenom,
				BurnRatio: "0.1",
			},
			setup: func(f *fixture) {
				burnAmount := sdk.NewInt64Coin(feeDenom, 100)
				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(burnAmount)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(burnAmount)).Return(nil),
				)
			},
			expectError: false,
		},
		{
			name: "successful burn with swap",
			fees: fees,
			params: types.Params{
				FeeDenom:   feeDenom,
				BurnDenom:  burnDenom,
				BurnRatio:  "0.1",
				BurnPoolId: 1,
			},
			setup: func(f *fixture) {
				amountToSwap := sdk.NewInt64Coin(feeDenom, 100)
				swappedAmount := sdk.NewInt64Coin(burnDenom, 95) // 5 is interface fee
				route := swaptypes.Route{
					DenomIn:  feeDenom,
					DenomOut: burnDenom,
					Strategy: &swaptypes.Route_Pool{Pool: &swaptypes.RoutePool{PoolId: 1}},
				}

				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(amountToSwap)).Return(nil),
					f.mocks.SwapKeeper.EXPECT().SwapExactAmountIn(
						gomock.Any(),
						authtypes.NewModuleAddress(types.ModuleName),
						authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
						route,
						amountToSwap.Amount,
						math.OneInt(),
					).Return(swaptypes.RouteResult{TokenOut: swappedAmount}, math.NewInt(5), nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(sdk.NewCoin(burnDenom, math.NewInt(90)))).Return(nil),
				)
			},
			expectError: false,
		},
		{
			name:        "no fee denom in fees",
			fees:        sdk.NewCoins(sdk.NewInt64Coin("other", 1000)),
			params:      types.DefaultParams(),
			setup:       func(f *fixture) {},
			expectError: false,
		},
		{
			name: "burn ratio is zero",
			fees: fees,
			params: types.Params{
				FeeDenom:  feeDenom,
				BurnRatio: "0.0",
			},
			setup:       func(f *fixture) {},
			expectError: false,
		},
		{
			name: "send coins fails",
			fees: fees,
			params: types.Params{
				FeeDenom:  feeDenom,
				BurnRatio: "0.1",
			},
			setup: func(f *fixture) {
				f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send error"))
			},
			expectError: false, // error is logged but not returned
		},
		{
			name: "swap fails",
			fees: fees,
			params: types.Params{
				FeeDenom:   feeDenom,
				BurnDenom:  burnDenom,
				BurnRatio:  "0.1",
				BurnPoolId: 1,
			},
			setup: func(f *fixture) {
				amountToSwap := sdk.NewInt64Coin(feeDenom, 100)
				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), authtypes.FeeCollectorName, types.ModuleName, sdk.NewCoins(amountToSwap)).Return(nil),
					f.mocks.SwapKeeper.EXPECT().SwapExactAmountIn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(swaptypes.RouteResult{}, math.ZeroInt(), errors.New("swap error")),
				)
			},
			expectError: false, // error is logged but not returned
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			sdkCtx := f.ctx.(sdk.Context)
			require.NoError(t, f.keeper.Params.Set(sdkCtx, tc.params))
			tc.setup(f)

			err := f.keeper.Burn(sdkCtx, tc.fees)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
