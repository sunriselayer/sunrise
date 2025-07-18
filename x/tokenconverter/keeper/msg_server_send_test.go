package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/tokenconverter/keeper"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

func TestMsgServer_Send(t *testing.T) {
	fromAddr := sdk.AccAddress("from")
	toAddr := sdk.AccAddress("to")
	allowedAddr := sdk.AccAddress("allowed")
	unrelatedAddr := sdk.AccAddress("unrelated")
	coin := sdk.NewCoins(sdk.NewCoin("test", math.NewInt(100)))
	zeroCoin := sdk.Coins{sdk.Coin{Denom: "test", Amount: math.NewInt(0)}}

	testCases := []struct {
		name        string
		msg         *types.MsgSend
		mockSetup   func(f *fixture)
		expectedErr string
	}{
		{
			name: "success - from address is allowed",
			msg:  &types.MsgSend{FromAddress: allowedAddr.String(), ToAddress: unrelatedAddr.String(), Amount: coin},
			mockSetup: func(f *fixture) {
				params := types.NewParams("", "", []string{allowedAddr.String()})
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), allowedAddr, unrelatedAddr, coin).Return(nil)
			},
		},
		{
			name: "success - to address is allowed",
			msg:  &types.MsgSend{FromAddress: fromAddr.String(), ToAddress: allowedAddr.String(), Amount: coin},
			mockSetup: func(f *fixture) {
				params := types.NewParams("", "", []string{allowedAddr.String()})
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), fromAddr, allowedAddr, coin).Return(nil)
			},
		},
		{
			name: "fail - from and to addresses are not allowed",
			msg:  &types.MsgSend{FromAddress: fromAddr.String(), ToAddress: toAddr.String(), Amount: coin},
			mockSetup: func(f *fixture) {
				params := types.NewParams("", "", []string{allowedAddr.String()})
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)
			},
			expectedErr: sdkerrors.ErrUnauthorized.Error(),
		},
		{
			name:        "fail - invalid from address",
			msg:         &types.MsgSend{FromAddress: "invalid", ToAddress: toAddr.String(), Amount: coin},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid from address",
		},
		{
			name:        "fail - invalid to address",
			msg:         &types.MsgSend{FromAddress: fromAddr.String(), ToAddress: "invalid", Amount: coin},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid to address",
		},
		{
			name:        "fail - zero amount",
			msg:         &types.MsgSend{FromAddress: fromAddr.String(), ToAddress: toAddr.String(), Amount: zeroCoin},
			mockSetup:   func(f *fixture) {},
			expectedErr: "amount must be positive",
		},
		{
			name: "fail - bank send error",
			msg:  &types.MsgSend{FromAddress: allowedAddr.String(), ToAddress: toAddr.String(), Amount: coin},
			mockSetup: func(f *fixture) {
				params := types.NewParams("", "", []string{allowedAddr.String()})
				err := f.keeper.SetParams(f.ctx, params)
				require.NoError(t, err)

				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), allowedAddr, toAddr, coin).Return(errors.New("bank send error"))
			},
			expectedErr: "bank send error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			msgServer := keeper.NewMsgServerImpl(f.keeper)
			_, err := msgServer.Send(f.ctx, tc.msg)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
