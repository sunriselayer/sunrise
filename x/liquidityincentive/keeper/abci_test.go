package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	stakingtypes "cosmossdk.io/x/staking/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

func TestCreateEpoch(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(tallyFixture)
		expectedTally []types.PoolWeight
		expectError   bool
	}{
		{
			name: "no votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 0)
			},
			expectedTally: []types.PoolWeight{},
		},
		{
			name: "one validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.PoolWeight{{PoolId: 1, Weight: "1000000"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := initFixture(t)
			ctx := fixture.ctx
			k := fixture.keeper
			mocks := getMocks(t)

			var (
				numVals       = 10
				numDelegators = 5
				addrs         = simtestutil.CreateRandomAccounts(numVals + numDelegators)
				valAddrs      = simtestutil.ConvertAddrsToValAddrs(addrs[:numVals])
				delAddrs      = addrs[numVals:]
			)
			// Mocks a bunch of validators
			mocks.StakingKeeper.EXPECT().
				IterateBondedValidatorsByPower(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, fn func(index int64, validator stakingtypes.Validator) bool) error {
						for i := int64(0); i < int64(numVals); i++ {
							valAddr, err := mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(valAddrs[i])
							require.NoError(t, err)
							fn(i, stakingtypes.Validator{
								OperatorAddress: valAddr,
								Status:          stakingtypes.Bonded,
								Tokens:          sdkmath.NewInt(1000000),
								DelegatorShares: sdkmath.LegacyNewDec(1000000),
							})
						}
						return nil
					})

			suite := tallyFixture{
				t:        t,
				valAddrs: valAddrs,
				delAddrs: delAddrs,
				ctx:      sdk.UnwrapSDKContext(ctx),
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			err := k.CreateEpoch(sdk.UnwrapSDKContext(ctx), 0, 1)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if len(tt.expectedTally) > 0 {
				epochs, err := k.GetAllEpoch(ctx)
				require.NoError(t, err)
				require.Len(t, epochs, 1)
				require.Equal(t, epochs[0].Id, uint64(1))
				require.Equal(t, epochs[0].StartBlock, int64(0))
				require.Equal(t, epochs[0].EndBlock, int64(5))
				require.Len(t, epochs[0].Gauges, 1)

				gauges := k.GetAllGauges(ctx)
				require.Len(t, gauges, 1)
				require.Equal(t, gauges[0].PreviousEpochId, uint64(0))
				require.Equal(t, gauges[0].PoolId, tt.expectedTally[0].PoolId)
				require.Equal(t, gauges[0].Count, tt.expectedTally[0].Weight)
			}
		})
	}
}

func TestEndBlocker(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(tallyFixture)
		expectedTally []types.PoolWeight
		expectError   bool
	}{
		{
			name: "no votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 0)
			},
			expectedTally: []types.PoolWeight{},
		},
		{
			name: "one validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.PoolWeight{{PoolId: 1, Weight: "1000000"}},
		},
		{
			name: "historical epochs",
			setup: func(s tallyFixture) {
				s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges:     []types.Gauge{},
				})
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.PoolWeight{{PoolId: 1, Weight: "1000000"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := initFixture(t)
			ctx := fixture.ctx
			k := fixture.keeper
			mocks := getMocks(t)

			var (
				numVals       = 10
				numDelegators = 5
				addrs         = simtestutil.CreateRandomAccounts(numVals + numDelegators)
				valAddrs      = simtestutil.ConvertAddrsToValAddrs(addrs[:numVals])
				delAddrs      = addrs[numVals:]
			)
			// Mocks a bunch of validators
			mocks.StakingKeeper.EXPECT().
				IterateBondedValidatorsByPower(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, fn func(index int64, validator stakingtypes.Validator) bool) error {
						for i := int64(0); i < int64(numVals); i++ {
							valAddr, err := mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(valAddrs[i])
							require.NoError(t, err)
							fn(i, stakingtypes.Validator{
								OperatorAddress: valAddr,
								Status:          stakingtypes.Bonded,
								Tokens:          sdkmath.NewInt(1000000),
								DelegatorShares: sdkmath.LegacyNewDec(1000000),
							})
						}
						return nil
					})

			suite := tallyFixture{
				t:        t,
				valAddrs: valAddrs,
				delAddrs: delAddrs,
				ctx:      sdk.UnwrapSDKContext(ctx),
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			err := k.EndBlocker(sdk.UnwrapSDKContext(ctx))
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if len(tt.expectedTally) > 0 {
				epochs, err := k.GetAllEpoch(ctx)
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(epochs), 1)
				epoch, found, err := k.GetLastEpoch(ctx)
				require.NoError(t, err)
				require.True(t, found)
				require.GreaterOrEqual(t, epoch.Id, uint64(1))
				require.Equal(t, epoch.StartBlock, int64(0))
				require.Equal(t, epoch.EndBlock, int64(5))
				require.Len(t, epoch.Gauges, 1)

				gauges := k.GetAllGauges(ctx)
				require.Len(t, gauges, 1)
				require.GreaterOrEqual(t, gauges[0].PreviousEpochId, uint64(0))
				require.Equal(t, gauges[0].PoolId, tt.expectedTally[0].PoolId)
				require.Equal(t, gauges[0].Count, tt.expectedTally[0].Weight)
			}
		})
	}
}

func TestBeginBlocker(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(tallyFixture)
		expectError bool
	}{
		{
			name: "empty epochs",
			setup: func(s tallyFixture) {
				s.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{sdk.NewInt64Coin(consts.BondDenom, 1000000)}).AnyTimes()
			},
		},
		{
			name: "existing epochs with positive fee collector balance",
			setup: func(s tallyFixture) {
				s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges: []types.Gauge{
						{
							PreviousEpochId: 0,
							PoolId:          1,
							Count:           math.OneInt(),
						},
					},
				})

				params, err := s.keeper.Params.Get(s.ctx)
				require.NoError(t, err)
				params.StakingRewardRatio = math.LegacyZeroDec().String()
				err = s.keeper.Params.Set(s.ctx, params)
				require.NoError(t, err)
				s.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{sdk.NewInt64Coin(consts.BondDenom, 1000000)}).AnyTimes()
				s.mocks.LiquiditypoolKeeper.EXPECT().AllocateIncentive(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).AnyTimes()
			},
		},
		{
			name: "zero liquidity incentive allocation",
			setup: func(s tallyFixture) {
				s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges: []types.Gauge{
						{
							PreviousEpochId: 0,
							PoolId:          1,
							Count:           math.OneInt(),
						},
					},
				})

				params, err := s.keeper.Params.Get(s.ctx)
				require.NoError(t, err)
				params.StakingRewardRatio = math.LegacyOneDec().String()
				err = s.keeper.Params.Set(s.ctx, params)
				require.NoError(t, err)
				s.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).
					Return(sdk.Coins{sdk.NewInt64Coin(consts.BondDenom, 1000000)}).AnyTimes()
				s.mocks.LiquiditypoolKeeper.EXPECT().AllocateIncentive(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).AnyTimes()
			},
		},
		{
			name: "existing epochs with empty fee collector balance",
			setup: func(s tallyFixture) {
				s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges: []types.Gauge{
						{
							PreviousEpochId: 0,
							PoolId:          1,
							Count:           math.OneInt(),
						},
					},
				})

				s.mocks.BankKeeper.EXPECT().GetAllBalances(gomock.Any(), gomock.Any()).Return(sdk.Coins{}).AnyTimes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := initFixture(t)
			ctx := fixture.ctx
			k := fixture.keeper
			mocks := getMocks(t)

			var (
				numVals       = 10
				numDelegators = 5
				addrs         = simtestutil.CreateRandomAccounts(numVals + numDelegators)
				valAddrs      = simtestutil.ConvertAddrsToValAddrs(addrs[:numVals])
				delAddrs      = addrs[numVals:]
			)
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
