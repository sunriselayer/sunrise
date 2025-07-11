package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// Test for InitLockupAccount msg server.
// It includes test cases for success, account already exists, and invalid owner address.
func TestMsgServer_InitLockupAccount(t *testing.T) {
	sender := sdk.AccAddress("sender")
	stakeDenom := "stake"

	testCases := []struct {
		name         string
		msg          *types.MsgInitLockupAccount
		mockSetup    func(f *fixture)
		expectedErr  string
		postRunCheck func(t *testing.T, f *fixture, msg *types.MsgInitLockupAccount)
	}{
		{
			name: "success",
			msg: &types.MsgInitLockupAccount{
				Sender: sender.String(),
				Owner:  sender.String(),
				Amount: sdk.NewCoin(stakeDenom, math.NewInt(0)),
			},
			mockSetup: func(f *fixture) {
				gomock.InOrder(
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(stakeDenom, nil),
					f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			postRunCheck: func(t *testing.T, f *fixture, msg *types.MsgInitLockupAccount) {
				owner, err := f.addressCodec.StringToBytes(msg.Owner)
				require.NoError(t, err)

				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.Equal(t, uint64(1), lockup.Id)
				require.Equal(t, msg.Owner, lockup.Owner)
			},
		},
		{
			name: "fail - invalid owner address",
			msg: &types.MsgInitLockupAccount{
				Sender: sender.String(),
				Owner:  "invalid-address",
				Amount: sdk.NewCoin(stakeDenom, math.NewInt(0)),
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid owner address: decoding bech32 failed: invalid separator index -1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			msgServer := keeper.NewMsgServerImpl(f.keeper)
			_, err := msgServer.InitLockupAccount(f.ctx, tc.msg)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				if tc.postRunCheck != nil {
					tc.postRunCheck(t, f, tc.msg)
				}
			}
		})
	}
}

// Test for InitLockupAccount with multiple accounts.
// It ensures that multiple accounts for the same owner can be created with different IDs.
func TestMsgServer_InitLockupAccount_MultipleAccounts(t *testing.T) {
	f := initFixture(t)
	msgServer := keeper.NewMsgServerImpl(f.keeper)
	sender := sdk.AccAddress("sender")
	stakeDenom := "stake"

	gomock.InOrder(
		f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(stakeDenom, nil),
		f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(stakeDenom, nil),
		f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)

	// Create first account
	_, err := msgServer.InitLockupAccount(f.ctx, &types.MsgInitLockupAccount{Sender: sender.String(), Owner: sender.String(), Amount: sdk.NewCoin(stakeDenom, math.NewInt(0))})
	require.NoError(t, err)

	owner, err := f.addressCodec.StringToBytes(sender.String())
	require.NoError(t, err)

	// Verify first account
	lockup1, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
	require.NoError(t, err)
	require.Equal(t, uint64(1), lockup1.Id)
	require.Equal(t, sender.String(), lockup1.Owner)
	require.NotEqual(t, "", lockup1.Address)
	require.True(t, lockup1.OriginalLocking.IsZero())
	require.Empty(t, lockup1.UnbondEntries.Entries)

	// Create second account
	_, err = msgServer.InitLockupAccount(f.ctx, &types.MsgInitLockupAccount{Sender: sender.String(), Owner: sender.String(), Amount: sdk.NewCoin(stakeDenom, math.NewInt(0))})
	require.NoError(t, err)

	// Verify second account
	lockup2, err := f.keeper.GetLockupAccount(f.ctx, owner, 2)
	require.NoError(t, err)
	require.Equal(t, uint64(2), lockup2.Id)
	require.Equal(t, sender.String(), lockup2.Owner)
	require.NotEqual(t, "", lockup2.Address)
	require.NotEqual(t, lockup1.Address, lockup2.Address) // Addresses should be different
}
