package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/tokenconverter/keeper"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func TestMsgServer_Convert(t *testing.T) {
	sender := sdk.AccAddress("sender")
	nonTransferableDenom := "nontransfer"
	transferableDenom := "transfer"

	testCases := []struct {
		name        string
		msg         *types.MsgConvert
		mockSetup   func(f *fixture)
		expectedErr string
	}{
		{
			name: "success",
			msg: &types.MsgConvert{
				Sender: sender.String(),
				Amount: math.NewInt(100),
			},
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
			name: "fail - invalid sender address",
			msg: &types.MsgConvert{
				Sender: "invalid",
				Amount: math.NewInt(100),
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid sender address",
		},
		{
			name: "fail - non-positive amount",
			msg: &types.MsgConvert{
				Sender: sender.String(),
				Amount: math.NewInt(0),
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "amount must be positive",
		},
		{
			name: "fail - keeper.Convert fails",
			msg: &types.MsgConvert{
				Sender: sender.String(),
				Amount: math.NewInt(100),
			},
			mockSetup: func(f *fixture) {
				params := types.NewParams(nonTransferableDenom, transferableDenom)
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				nonTransferableCoin := sdk.NewCoin(nonTransferableDenom, math.NewInt(100))
				f.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), sender, types.ModuleName, sdk.NewCoins(nonTransferableCoin)).Return(errors.New("bank send error"))
			},
			expectedErr: "bank send error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			msgServer := keeper.NewMsgServerImpl(f.keeper)
			_, err := msgServer.Convert(f.ctx, tc.msg)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
