package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec/address"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func TestStartNewEpoch(t *testing.T) {
	var (
		addrs      = simtestutil.CreateRandomAccounts(1)
		moduleAddr = simtestutil.CreateRandomAccounts(1)[0]
		valAddrs   = simtestutil.ConvertAddrsToValAddrs(addrs)
	)

	tests := []struct {
		name        string
		setup       func(fx *fixture, k keeper.Keeper, mocks LiquidityIncentiveMocks)
		expectError bool
	}{
		{
			name: "no votes - should fail",
			setup: func(fx *fixture, k keeper.Keeper, mocks LiquidityIncentiveMocks) {
				mocks.StakingKeeper.EXPECT().IterateBondedValidatorsByPower(gomock.Any(), gomock.Any()).Return(nil)
				mocks.AcctKeeper.EXPECT().GetModuleAddress(shareclasstypes.ModuleName).Return(moduleAddr)
				mocks.StakingKeeper.EXPECT().IterateDelegations(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mocks.StakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(sdkmath.NewInt(0), nil)
			},
			expectError: true,
		},
		{
			name: "one validator votes - success",
			setup: func(fx *fixture, k keeper.Keeper, mocks LiquidityIncentiveMocks) {
				// Setup vote
				vote := types.Vote{Sender: sdk.AccAddress(valAddrs[0]).String(), PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "1"}}}
				err := k.SetVote(fx.ctx, vote)
				require.NoError(t, err)

				// Mocks
				mocks.StakingKeeper.EXPECT().IterateBondedValidatorsByPower(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) bool) error {
						fn(0, stakingtypes.Validator{OperatorAddress: valAddrs[0].String(), Status: stakingtypes.Bonded, Tokens: sdkmath.NewInt(1000000), DelegatorShares: sdkmath.LegacyNewDec(1000000)})
						return nil
					})
				mocks.AcctKeeper.EXPECT().GetModuleAddress(shareclasstypes.ModuleName).Return(moduleAddr)
				mocks.StakingKeeper.EXPECT().IterateDelegations(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mocks.StakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(sdkmath.NewInt(1000000), nil)
				mocks.AcctKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("cosmos"))
				mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("cosmosvaloper")).AnyTimes()
				// Mocks for FinalizeBribeForEpoch
				mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(sdk.AccAddress("fee_collector")).AnyTimes()
			},
			expectError: false,
		},
		{
			name: "epoch not ended - should fail",
			setup: func(fx *fixture, k keeper.Keeper, mocks LiquidityIncentiveMocks) {
				err := k.SetEpoch(fx.ctx, types.Epoch{
					Id:       1,
					EndBlock: sdk.UnwrapSDKContext(fx.ctx).BlockHeight() + 1,
				})
				require.NoError(t, err)
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fx := initFixture(t)
			msgServer := keeper.NewMsgServerImpl(fx.keeper)

			if tt.setup != nil {
				tt.setup(fx, fx.keeper, fx.mocks)
			}

			_, err := msgServer.StartNewEpoch(fx.ctx, &types.MsgStartNewEpoch{
				Sender: addrs[0].String(),
			})

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				epoch, found, err := fx.keeper.GetLastEpoch(fx.ctx)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, uint64(1), epoch.Id)
			}
		})
	}
}
