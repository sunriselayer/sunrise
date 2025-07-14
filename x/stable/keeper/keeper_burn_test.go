// keeper_burn_test.go
//
// This file contains the test suite for the burning functionality of the keeper.
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

func TestKeeperBurn(t *testing.T) {
	authority := sdk.AccAddress("authority")
	sender := sdk.AccAddress("sender")

	tests := []struct {
		name        string
		setup       func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context)
		sender      sdk.AccAddress
		amount      sdk.Coins
		expectedErr *errors.Error
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
			sender: authority,
			amount: sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(100))),
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
			sender:      sender,
			amount:      sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(100))),
			expectedErr: types.ErrInvalidAuthority,
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
			sender:      authority,
			amount:      sdk.NewCoins(sdk.NewCoin("stable", math.NewInt(0))),
			expectedErr: types.ErrInvalidAmount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, mocks, _, sdkCtx := setupMsgServer(t)
			if tc.setup != nil {
				tc.setup(k, mocks.BankKeeper, sdk.UnwrapSDKContext(sdkCtx))
			}

			err := k.Burn(sdkCtx, tc.sender, tc.amount)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
