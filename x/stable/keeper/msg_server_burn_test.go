// msg_server_burn_test.go
//
// This file contains the test suite for the burn functionality of the message server.
package keeper_test

import (
	"testing"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/stable/keeper"
	stabletestutil "github.com/sunriselayer/sunrise/x/stable/testutil"
	"github.com/sunriselayer/sunrise/x/stable/types"
)

func TestMsgServerBurn(t *testing.T) {
	authority := sdk.AccAddress("authority")
	sender := sdk.AccAddress("sender")

	tests := []struct {
		name           string
		setup          func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context)
		msg            *types.MsgBurn
		expectedErr    *errors.Error
		expectedErrStr string
	}{
		{
			name: "success: burn as authority",
			setup: func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context) {
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)

				params.AuthorityAddresses = []string{authority.String()}
				err = k.Params.Set(ctx, params)
				require.NoError(t, err)

				bankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), authority, types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().BurnCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
			},
			msg: &types.MsgBurn{
				Sender: authority.String(),
				Amount: sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(100))),
			},
		},
		{
			name: "failure: burn as non-authority",
			setup: func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context) {
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)

				params.AuthorityAddresses = []string{authority.String()}
				err = k.Params.Set(ctx, params)
				require.NoError(t, err)
			},
			msg: &types.MsgBurn{
				Sender: sender.String(),
				Amount: sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(100))),
			},
			expectedErr: types.ErrInvalidAuthority,
		},
		{
			name: "failure: invalid sender address",
			msg: &types.MsgBurn{
				Sender: "invalid-address",
				Amount: sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(100))),
			},
			expectedErrStr: "invalid sender address",
		},
		{
			name: "failure: non-positive burn amount",
			setup: func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context) {
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)

				params.AuthorityAddresses = []string{authority.String()}
				err = k.Params.Set(ctx, params)
				require.NoError(t, err)
			},
			msg: &types.MsgBurn{
				Sender: authority.String(),
				Amount: sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(0))),
			},
			expectedErr: types.ErrInvalidAmount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, mocks, srv, sdkCtx := setupMsgServer(t)
			if tc.setup != nil {
				tc.setup(k, mocks.BankKeeper, sdk.UnwrapSDKContext(sdkCtx))
			}

			_, err := srv.Burn(sdkCtx, tc.msg)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else if tc.expectedErrStr != "" {
				require.ErrorContains(t, err, tc.expectedErrStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
