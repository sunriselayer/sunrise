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

// Test for NonVotingUndelegate msg server, based on the structure of TestMsgServer_NonVotingDelegate.
func TestMsgServer_NonVotingUndelegate(t *testing.T) {
	owner := sdk.AccAddress("owner")
	valAddr := sdk.ValAddress("validator")
	transferableDenom := "transferable"
	valAddressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	completionTime := time.Now().Add(time.Hour)

	testCases := []struct {
		name         string
		msg          *types.MsgNonVotingUndelegate
		mockSetup    func(f *fixture)
		expectedErr  string
		postRunCheck func(t *testing.T, f *fixture)
	}{
		{
			name: "success - new unbonding entry",
			msg: &types.MsgNonVotingUndelegate{
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
					f.mocks.ShareclassKeeper.EXPECT().Undelegate(gomock.Any(), lockupAddr, lockupAddr, valAddr, gomock.Any()).
						Return(sdk.NewCoin(transferableDenom, math.NewInt(100)), sdk.NewCoins(), completionTime, nil),
				)
			},
			postRunCheck: func(t *testing.T, f *fixture) {
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.Len(t, lockup.UnbondEntries.Entries, 1)
				require.True(t, lockup.UnbondEntries.Entries[0].Amount.Equal(math.NewInt(100)))
			},
		},
		{
			name: "success - merge with existing unbonding entry",
			msg: &types.MsgNonVotingUndelegate{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewCoin(transferableDenom, math.NewInt(50)),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				sdkCtx := sdk.UnwrapSDKContext(f.ctx)
				lockup := types.LockupAccount{
					Owner:   owner.String(),
					Id:      1,
					Address: lockupAddr.String(),
					UnbondEntries: &types.UnbondingEntries{
						Entries: []*types.UnbondingEntry{
							{CreationHeight: sdkCtx.BlockHeight(), EndTime: completionTime.Unix(), Amount: math.NewInt(100)},
						},
					},
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				gomock.InOrder(
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil),
					f.mocks.ShareclassKeeper.EXPECT().Undelegate(gomock.Any(), lockupAddr, lockupAddr, valAddr, gomock.Any()).
						Return(sdk.NewCoin(transferableDenom, math.NewInt(50)), sdk.NewCoins(), completionTime, nil),
				)
			},
			postRunCheck: func(t *testing.T, f *fixture) {
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.Len(t, lockup.UnbondEntries.Entries, 1)
				require.True(t, lockup.UnbondEntries.Entries[0].Amount.Equal(math.NewInt(150))) // 100 + 50
			},
		},
		{
			name: "fail - invalid owner address",
			msg: &types.MsgNonVotingUndelegate{
				Owner: "invalid",
			},
			mockSetup:   func(f *fixture) {},
			expectedErr: "invalid owner address",
		},
		{
			name: "fail - invalid validator address",
			msg: &types.MsgNonVotingUndelegate{
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
			msg: &types.MsgNonVotingUndelegate{
				Owner:            owner.String(),
				LockupAccountId:  99,
				ValidatorAddress: valAddr.String(),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
			},
			expectedErr: "lockup account not found",
		},
		{
			name: "fail - invalid denom",
			msg: &types.MsgNonVotingUndelegate{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewCoin("wrongdenom", math.NewInt(100)),
			},
			mockSetup: func(f *fixture) {
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{Owner: owner.String(), Id: 1, Address: lockupAddr.String()}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil)
			},
			expectedErr: "invalid denom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			msgServer := keeper.NewMsgServerImpl(f.keeper)
			_, err := msgServer.NonVotingUndelegate(f.ctx, tc.msg)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
				if tc.postRunCheck != nil {
					tc.postRunCheck(t, f)
				}
			}
		})
	}
}
