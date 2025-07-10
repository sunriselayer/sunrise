// This file is a unit test for the keeper's lockup account logic.
// It covers the following functions:
// - LockupAccountAddress
// - GetAndIncrementNextLockupAccountID
// - InitLockupAccountFromMsg
package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func TestKeeper_LockupAccountAddress(t *testing.T) {
	f := initFixture(t)
	owner := sdk.AccAddress("owner")
	id1 := uint64(1)
	id2 := uint64(2)

	// Check for determinism
	addr1_a := f.keeper.LockupAccountAddress(owner, id1)
	addr1_b := f.keeper.LockupAccountAddress(owner, id1)
	require.Equal(t, addr1_a, addr1_b)

	// Check for uniqueness
	addr2 := f.keeper.LockupAccountAddress(owner, id2)
	require.NotEqual(t, addr1_a, addr2)
}

func TestKeeper_GetAndIncrementNextLockupAccountID(t *testing.T) {
	f := initFixture(t)
	owner := sdk.AccAddress("owner")

	// First time for an owner
	currentID, nextID, err := f.keeper.GetAndIncrementNextLockupAccountID(f.ctx, owner)
	require.NoError(t, err)
	require.Equal(t, uint64(1), currentID)
	require.Equal(t, uint64(2), nextID)

	// Second time for the same owner
	currentID, nextID, err = f.keeper.GetAndIncrementNextLockupAccountID(f.ctx, owner)
	require.NoError(t, err)
	require.Equal(t, uint64(2), currentID)
	require.Equal(t, uint64(3), nextID)
}

func TestKeeper_InitLockupAccountFromMsg(t *testing.T) {
	sender := sdk.AccAddress("sender")
	owner := sdk.AccAddress("owner")
	transferableDenom := "transferable"

	testCases := []struct {
		name        string
		msg         *types.MsgInitLockupAccount
		mockSetup   func(f *fixture)
		expectedErr string
	}{
		{
			name: "success",
			msg: &types.MsgInitLockupAccount{
				Sender:    sender.String(),
				Owner:     owner.String(),
				Amount:    sdk.NewCoin(transferableDenom, math.NewInt(100)),
				StartTime: time.Now().Unix(),
				EndTime:   time.Now().Add(time.Hour).Unix(),
			},
			mockSetup: func(f *fixture) {
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "fail - invalid sender",
			msg: &types.MsgInitLockupAccount{
				Sender: "invalid",
				Owner:  owner.String(),
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid sender address",
		},
		{
			name: "fail - invalid owner",
			msg: &types.MsgInitLockupAccount{
				Sender: sender.String(),
				Owner:  "invalid",
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid owner address",
		},
		{
			name: "fail - invalid denom",
			msg: &types.MsgInitLockupAccount{
				Sender: sender.String(),
				Owner:  owner.String(),
				Amount: sdk.NewCoin("wrongdenom", math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
			},
			expectedErr: "invalid denom",
		},
		{
			name: "fail - send coins failed",
			msg: &types.MsgInitLockupAccount{
				Sender: sender.String(),
				Owner:  owner.String(),
				Amount: sdk.NewCoin(transferableDenom, math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(types.ErrInsufficientUnlockedFunds)
			},
			expectedErr: "failed to send coins",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			err := f.keeper.InitLockupAccountFromMsg(f.ctx, tc.msg)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				// Verify account was created with ID 1
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Owner, lockup.Owner)
				require.Equal(t, tc.msg.Amount.Amount, lockup.OriginalLocking)

				// Verify next ID is 2
				nextID, err := f.keeper.NextLockupAccountId.Get(f.ctx, owner)
				require.NoError(t, err)
				require.Equal(t, uint64(2), nextID)
			}
		})
	}
}
