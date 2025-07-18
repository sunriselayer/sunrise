// keeper_mint_test.go
//
// This file contains the test suite for the minting functionality of the keeper.
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

func TestKeeperMint(t *testing.T) {
	authority := sdk.AccAddress("authority")
	sender := sdk.AccAddress("sender")

	tests := []struct {
		name        string
		setup       func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context)
		sender      sdk.AccAddress
		amount      math.Int
		expectedErr *errors.Error
	}{
		{
			name: "success: mint as authority",
			setup: func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context) {
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)

				params.AllowedAddresses = []string{authority.String()}
				err = k.Params.Set(ctx, params)
				require.NoError(t, err)

				bankKeeper.EXPECT().MintCoins(gomock.Any(), types.ModuleName, gomock.Any()).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, authority, gomock.Any()).Return(nil)
			},
			sender: authority,
			amount: math.NewInt(100),
		},
		{
			name: "failure: mint as non-authority",
			setup: func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context) {
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)

				params.AllowedAddresses = []string{authority.String()}
				err = k.Params.Set(ctx, params)
				require.NoError(t, err)
			},
			sender:      sender,
			amount:      math.NewInt(100),
			expectedErr: types.ErrInvalidAuthority,
		},
		{
			name: "failure: non-positive mint amount",
			setup: func(k keeper.Keeper, bankKeeper *stabletestutil.MockBankKeeper, ctx sdk.Context) {
				params, err := k.Params.Get(ctx)
				require.NoError(t, err)

				params.AllowedAddresses = []string{authority.String()}
				err = k.Params.Set(ctx, params)
				require.NoError(t, err)
			},
			sender:      authority,
			amount:      math.NewInt(0),
			expectedErr: types.ErrInvalidAmount,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, mocks, _, sdkCtx := setupMsgServer(t)
			if tc.setup != nil {
				tc.setup(k, mocks.BankKeeper, sdk.UnwrapSDKContext(sdkCtx))
			}

			_, err := k.Mint(sdkCtx, tc.sender, tc.amount)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
