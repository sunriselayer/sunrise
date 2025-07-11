// This file is a unit test for the Send function in the msg_server_send.go file.
package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

func TestMsgServer_Send(t *testing.T) {
	// Setup
	f := initFixture(t)
	msgServer := keeper.NewMsgServerImpl(f.keeper)
	owner := sdk.AccAddress("owner")
	recipient := sdk.AccAddress("recipient")
	transferableDenom := "stake"

	// Create a lockup account
	lockupAccount := types.LockupAccount{
		Owner:           owner.String(),
		Id:              0,
		Address:         sdk.AccAddress("lockup_acct").String(),
		OriginalLocking: math.NewInt(1000),
		UnbondEntries:   &types.UnbondingEntries{},
	}
	ownerAddr, err := sdk.AccAddressFromBech32(lockupAccount.Owner)
	require.NoError(t, err)
	pair := collections.Join(ownerAddr.Bytes(), lockupAccount.Id)
	err = f.keeper.LockupAccounts.Set(f.ctx, pair, lockupAccount)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		msg           *types.MsgSend
		mockSetup     func(f *fixture)
		expectErr     bool
		expectedError string
	}{
		{
			name: "successful send",
			msg: &types.MsgSend{
				Owner:           owner.String(),
				LockupAccountId: 0,
				Recipient:       recipient.String(),
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100)),
			},
			mockSetup: func(f *fixture) {
				// Mock dependencies
				recipientAddr, err := sdk.AccAddressFromBech32(recipient.String())
				require.NoError(t, err)

				lockupAddr, err := sdk.AccAddressFromBech32(lockupAccount.Address)
				require.NoError(t, err)

				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), lockupAddr, transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(1000)))
				// Since CheckUnbondingEntriesMature does not use staking keeper, we don't need to mock it.
				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), lockupAddr, recipientAddr, sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100))).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "successful send with multiple denoms",
			msg: &types.MsgSend{
				Owner:           owner.String(),
				LockupAccountId: 0,
				Recipient:       recipient.String(),
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100), sdk.NewInt64Coin("otherdenom", 50)),
			},
			mockSetup: func(f *fixture) {
				recipientAddr, err := sdk.AccAddressFromBech32(recipient.String())
				require.NoError(t, err)
				lockupAddr, err := sdk.AccAddressFromBech32(lockupAccount.Address)
				require.NoError(t, err)

				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				// Only transferableDenom's balance is checked
				f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), lockupAddr, transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(1000)))
				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), lockupAddr, recipientAddr, sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100), sdk.NewInt64Coin("otherdenom", 50))).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "sendable balance is smaller",
			msg: &types.MsgSend{
				Owner:           owner.String(),
				LockupAccountId: 0,
				Recipient:       recipient.String(),
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 1001)),
			},
			mockSetup: func(f *fixture) {
				lockupAddr, err := sdk.AccAddressFromBech32(lockupAccount.Address)
				require.NoError(t, err)

				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), lockupAddr, transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(1000)))
			},
			expectErr:     true,
			expectedError: "spendable balance 1000stake is smaller than 1001stake: not enough spendable balance",
		},
		{
			name: "bank keeper send coins failed",
			msg: &types.MsgSend{
				Owner:           owner.String(),
				LockupAccountId: 0,
				Recipient:       recipient.String(),
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100)),
			},
			mockSetup: func(f *fixture) {
				recipientAddr, err := sdk.AccAddressFromBech32(recipient.String())
				require.NoError(t, err)
				lockupAddr, err := sdk.AccAddressFromBech32(lockupAccount.Address)
				require.NoError(t, err)

				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
				f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), lockupAddr, transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(1000)))
				f.mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), lockupAddr, recipientAddr, sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100))).Return(sdkerrors.ErrInsufficientFunds)
			},
			expectErr:     true,
			expectedError: "insufficient funds",
		},
		{
			name: "invalid owner address",
			msg: &types.MsgSend{
				Owner:           "invalid_address",
				LockupAccountId: 0,
				Recipient:       recipient.String(),
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
			},
			expectErr:     true,
			expectedError: "invalid owner address",
		},
		{
			name: "invalid recipient address",
			msg: &types.MsgSend{
				Owner:           owner.String(),
				LockupAccountId: 0,
				Recipient:       "invalid_address",
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
			},
			expectErr:     true,
			expectedError: "invalid recipient address",
		},
		{
			name: "lockup account not found",
			msg: &types.MsgSend{
				Owner:           owner.String(),
				LockupAccountId: 1, // Non-existent account ID
				Recipient:       recipient.String(),
				Amount:          sdk.NewCoins(sdk.NewInt64Coin(transferableDenom, 100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.AccountKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
			},
			expectErr:     true,
			expectedError: "lockup account not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Redo setup for each test case to ensure isolation
			f := initFixture(t)
			msgServer = keeper.NewMsgServerImpl(f.keeper)
			ownerAddr, err := sdk.AccAddressFromBech32(lockupAccount.Owner)
			require.NoError(t, err)
			pair := collections.Join(ownerAddr.Bytes(), lockupAccount.Id)
			// Reset the account for each test, in case a previous test case removed it.
			err = f.keeper.LockupAccounts.Set(f.ctx, pair, lockupAccount)
			require.NoError(t, err)

			sdkCtx := f.ctx.(sdk.Context).WithBlockTime(time.Now())

			if tc.mockSetup != nil {
				tc.mockSetup(f)
			}

			// Execute
			_, err = msgServer.Send(sdk.WrapSDKContext(sdkCtx), tc.msg)

			// Assert
			if tc.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
