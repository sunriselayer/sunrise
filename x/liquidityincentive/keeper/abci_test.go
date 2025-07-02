package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestBeginBlocker(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(tallyFixture)
		expectError bool
	}{
		{
			name: "empty epochs",
			setup: func(s tallyFixture) {
				s.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.NewInt64Coin(consts.BondDenom, 1000000)).AnyTimes()
			},
		},
		{
			name: "existing epochs with positive fee collector balance",
			setup: func(s tallyFixture) {
				err := s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges: []types.Gauge{
						{
							PoolId:      1,
							VotingPower: math.OneInt(),
						},
					},
				})
				require.NoError(t, err)
				params, err := s.keeper.Params.Get(s.ctx)
				require.NoError(t, err)
				params.StakingRewardRatio = math.LegacyZeroDec().String()
				err = s.keeper.Params.Set(s.ctx, params)
				require.NoError(t, err)
				s.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.NewInt64Coin(consts.BondDenom, 1000000)).AnyTimes()
			},
		},
		{
			name: "zero liquidity incentive allocation",
			setup: func(s tallyFixture) {
				err := s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges: []types.Gauge{
						{
							PoolId:      1,
							VotingPower: math.OneInt(),
						},
					},
				})
				require.NoError(t, err)
				params, err := s.keeper.Params.Get(s.ctx)
				require.NoError(t, err)
				params.StakingRewardRatio = math.LegacyOneDec().String()
				err = s.keeper.Params.Set(s.ctx, params)
				require.NoError(t, err)
				s.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.NewInt64Coin(consts.BondDenom, 1000000)).AnyTimes()
			},
		},
		{
			name: "existing epochs with empty fee collector balance",
			setup: func(s tallyFixture) {
				err := s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges: []types.Gauge{
						{
							PoolId:      1,
							VotingPower: math.OneInt(),
						},
					},
				})
				require.NoError(t, err)
				s.mocks.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sdk.NewInt64Coin(consts.BondDenom, 0)).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := f.ctx
			k := f.keeper
			mocks := f.mocks

			var (
				numVals       = 10
				numDelegators = 5
				addrs         = simtestutil.CreateRandomAccounts(numVals + numDelegators)
				valAddrs      = simtestutil.ConvertAddrsToValAddrs(addrs[:numVals])
				delAddrs      = addrs[numVals:]
			)
			mocks.StakingKeeper.EXPECT().BondDenom(gomock.Any()).
				Return(consts.BondDenom, nil).AnyTimes()
			mocks.FeeKeeper.EXPECT().FeeDenom(gomock.Any()).
				Return(consts.StableDenom, nil).AnyTimes()
			mocks.TokenConverterKeeper.EXPECT().ConvertReverse(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			mocks.BankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			mocks.TokenConverterKeeper.EXPECT().ConvertReverse(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			mocks.LiquiditypoolKeeper.EXPECT().AllocateIncentive(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			suite := tallyFixture{
				t:        t,
				valAddrs: valAddrs,
				delAddrs: delAddrs,
				ctx:      sdk.UnwrapSDKContext(ctx),
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			err := k.BeginBlocker(sdk.UnwrapSDKContext(ctx))
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
