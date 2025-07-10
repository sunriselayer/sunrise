package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// Test for NonVotingDelegate msg server.
// Covers success, invalid addresses, account not found, invalid denom, and delegation failure.
func TestMsgServer_NonVotingDelegate(t *testing.T) {
	owner := sdk.AccAddress("owner")
	valAddr := sdk.ValAddress("validator")
	transferableDenom := "transferable"
	valAddressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

	testCases := []struct {
		name        string
		msg         *types.MsgNonVotingDelegate
		mockSetup   func(f *fixture)
		expectedErr string
	}{
		{
			name: "success",
			msg: &types.MsgNonVotingDelegate{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewCoin(transferableDenom, math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)

				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{
					Owner:           owner.String(),
					Id:              1,
					Address:         lockupAddr.String(),
					OriginalLocking: math.NewInt(1000),
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				// Setup mocks for a successful delegation
				gomock.InOrder(
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil),
					f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(1000))),
					f.mocks.ShareclassKeeper.EXPECT().Delegate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sdk.NewCoins(sdk.NewCoin(transferableDenom, math.NewInt(100))), sdk.NewCoins(sdk.NewCoin(transferableDenom, math.NewInt(5))), nil),
				)
			},
		},
		{
			name: "fail - invalid owner address",
			msg: &types.MsgNonVotingDelegate{
				Owner: "invalid",
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid owner address",
		},
		{
			name: "fail - invalid validator address",
			msg: &types.MsgNonVotingDelegate{
				Owner:            owner.String(),
				ValidatorAddress: "invalid",
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
			},
			expectedErr: "invalid validator address",
		},
		{
			name: "fail - lockup account not found",
			msg: &types.MsgNonVotingDelegate{
				Owner:            owner.String(),
				LockupAccountId:  99,
				ValidatorAddress: valAddr.String(),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
			},
			expectedErr: types.ErrLockupAccountNotFound.Error(),
		},
		{
			name: "fail - invalid denom",
			msg: &types.MsgNonVotingDelegate{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewCoin("wrongdenom", math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{
					Owner:   owner.String(),
					Id:      1,
					Address: lockupAddr.String(),
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
			},
			expectedErr: "invalid denom",
		},
		{
			name: "fail - insufficient total balance",
			msg: &types.MsgNonVotingDelegate{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewCoin(transferableDenom, math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{
					Owner:   owner.String(),
					Id:      1,
					Address: lockupAddr.String(),
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				gomock.InOrder(
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil),
					f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(50))),
				)
			},
			expectedErr: "total balance is less than delegation amount",
		},
		{
			name: "fail - insufficient unlocked funds for delegation",
			msg: &types.MsgNonVotingDelegate{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewCoin(transferableDenom, math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				// Set a specific time for the context
				now := time.Now()
				sdkCtx := sdk.UnwrapSDKContext(f.ctx)
				f.ctx = sdkCtx.WithBlockTime(now)

				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				// Setup a lockup account with less unlocked funds than delegation amount
				lockup := types.LockupAccount{
					Owner:            owner.String(),
					Id:               1,
					Address:          lockupAddr.String(),
					OriginalLocking:  math.NewInt(500), // Total locked
					DelegatedLocking: math.NewInt(450), // Already delegated
					StartTime:        now.Unix(),
					EndTime:          now.Add(time.Hour).Unix(),
				}
				// Available for delegation = 50, but trying to delegate 100
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				gomock.InOrder(
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil),
					f.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), transferableDenom).Return(sdk.NewCoin(transferableDenom, math.NewInt(500))),
				)
			},
			expectedErr: types.ErrInsufficientUnlockedFunds.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			msgServer := keeper.NewMsgServerImpl(f.keeper)
			_, err := msgServer.NonVotingDelegate(f.ctx, tc.msg)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
