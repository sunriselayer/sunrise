// Tests for the Convert and ConvertReverse functions in the tokenconverter keeper.
// It covers success cases and various failure scenarios for both functions.
package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func TestKeeper_Convert(t *testing.T) {
	sender := sdk.AccAddress("sender")
	nonTransferableDenom := "nontransfer"
	transferableDenom := "transfer"

	testCases := []struct {
		name        string
		amount      math.Int
		mockSetup   func(f *fixture)
		expectedErr string
	}{
		{
			name:   "success",
			amount: math.NewInt(100),
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))

				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sender, sdk.NewCoins(transferableCoin)).Return(nil),
				)
			},
		},
		{
			name:        "fail - send coins to module fails",
			amount:      math.NewInt(100),
			expectedErr: "send coins failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(errors.New("send coins failed"))
			},
		},
		{
			name:        "fail - burn coins fails",
			amount:      math.NewInt(100),
			expectedErr: "burn coins failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(errors.New("burn coins failed")),
				)
			},
		},
		{
			name:        "fail - mint coins fails",
			amount:      math.NewInt(100),
			expectedErr: "mint coins failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(errors.New("mint coins failed")),
				)
			},
		},
		{
			name:        "fail - send coins from module fails",
			amount:      math.NewInt(100),
			expectedErr: "send coins from module failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))

				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sender, sdk.NewCoins(transferableCoin)).Return(errors.New("send coins from module failed")),
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			err := f.keeper.Convert(f.ctx, tc.amount, sender)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_ConvertReverse(t *testing.T) {
	sender := sdk.AccAddress("sender")
	nonTransferableDenom := "nontransfer"
	transferableDenom := "transfer"

	testCases := []struct {
		name        string
		amount      math.Int
		mockSetup   func(f *fixture)
		expectedErr string
	}{
		{
			name:   "success",
			amount: math.NewInt(100),
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))

				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sender, sdk.NewCoins(nonTransferableCoin)).Return(nil),
				)
			},
		},
		{
			name:        "fail - send coins to module fails",
			amount:      math.NewInt(100),
			expectedErr: "send coins failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))
				f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(transferableCoin)).Return(errors.New("send coins failed"))
			},
		},
		{
			name:        "fail - burn coins fails",
			amount:      math.NewInt(100),
			expectedErr: "burn coins failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))
				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(transferableCoin)).Return(errors.New("burn coins failed")),
				)
			},
		},
		{
			name:        "fail - mint coins fails",
			amount:      math.NewInt(100),
			expectedErr: "mint coins failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))
				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(errors.New("mint coins failed")),
				)
			},
		},
		{
			name:        "fail - send coins from module fails",
			amount:      math.NewInt(100),
			expectedErr: "send coins from module failed",
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				transferableCoin := sdk.NewCoin(transferableDenom, math.NewInt(100))

				gomock.InOrder(
					f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(transferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(nil),
					f.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, sender, sdk.NewCoins(nonTransferableCoin)).Return(errors.New("send coins from module failed")),
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			err := f.keeper.ConvertReverse(f.ctx, tc.amount, sender)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
