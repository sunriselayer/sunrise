package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/lockup/keeper"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// Test for ClaimRewards msg server.
// It covers success cases with and without transferable rewards, and various failure scenarios.
func TestMsgServer_ClaimRewards(t *testing.T) {
	owner := sdk.AccAddress(("owner"))
	valAddr := sdk.ValAddress([]byte("validator"))
	transferableDenom := "transferable"
	valAddressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

	testCases := []struct {
		name         string
		msg          *types.MsgClaimRewards
		mockSetup    func(f *fixture)
		expectedErr  string
		postRunCheck func(t *testing.T, f *fixture)
	}{
		{
			name: "success",
			msg: &types.MsgClaimRewards{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{
					Owner:             owner.String(),
					Id:                1,
					Address:           lockupAddr.String(),
					AdditionalLocking: math.ZeroInt(),
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				rewards := sdk.NewCoins(sdk.NewCoin(transferableDenom, math.NewInt(10)))
				gomock.InOrder(
					f.mocks.ShareclassKeeper.EXPECT().ClaimRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(rewards, nil),
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil),
				)
			},
			postRunCheck: func(t *testing.T, f *fixture) {
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				// check that rewards are added to AdditionalLocking
				require.Equal(t, math.NewInt(10), lockup.AdditionalLocking)
			},
		},
		{
			name: "success - no transferable rewards",
			msg: &types.MsgClaimRewards{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{
					Owner:             owner.String(),
					Id:                1,
					Address:           lockupAddr.String(),
					AdditionalLocking: math.ZeroInt(),
				}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				rewards := sdk.NewCoins(sdk.NewCoin("otherdenom", math.NewInt(10)))
				gomock.InOrder(
					f.mocks.ShareclassKeeper.EXPECT().ClaimRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(rewards, nil),
					f.mocks.TokenConverterKeeper.EXPECT().GetTransferableDenom(gomock.Any()).Return(transferableDenom, nil),
				)
			},
			postRunCheck: func(t *testing.T, f *fixture) {
				lockup, err := f.keeper.GetLockupAccount(f.ctx, owner, 1)
				require.NoError(t, err)
				require.True(t, lockup.AdditionalLocking.IsZero())
			},
		},
		{
			name: "fail - shareclass claim rewards fails",
			msg: &types.MsgClaimRewards{
				Owner:            owner.String(),
				LockupAccountId:  1,
				ValidatorAddress: valAddr.String(),
			},
			mockSetup: func(f *fixture) {
				f.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(valAddressCodec)
				lockupAddr := f.keeper.LockupAccountAddress(owner, 1)
				lockup := types.LockupAccount{Owner: owner.String(), Id: 1, Address: lockupAddr.String()}
				err := f.keeper.SetLockupAccount(f.ctx, lockup)
				require.NoError(t, err)

				gomock.InOrder(
					f.mocks.ShareclassKeeper.EXPECT().ClaimRewards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("claim failed")),
				)
			},
			expectedErr: "claim failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := initFixture(t)
			tc.mockSetup(f)

			msgServer := keeper.NewMsgServerImpl(f.keeper)
			_, err := msgServer.ClaimRewards(f.ctx, tc.msg)

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
