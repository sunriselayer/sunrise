package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sdkmath "cosmossdk.io/math"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func TestStartNewEpoch(t *testing.T) {
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
				err := s.keeper.SetEpoch(s.ctx, types.Epoch{
					Id:         1,
					StartBlock: 0,
					EndBlock:   0,
					Gauges:     []types.Gauge{},
				})
				require.NoError(t, err)
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.PoolWeight{{PoolId: 1, Weight: "1000000"}},
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
				moduleAddr    = simtestutil.CreateRandomAccounts(2)
			)
			// Mocks a bunch of validators
			mocks.StakingKeeper.EXPECT().
				IterateBondedValidatorsByPower(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) bool) error {
						for i := int64(0); i < int64(numVals); i++ {
							valAddr := valAddrs[i].String()
							fn(i, stakingtypes.Validator{
								OperatorAddress: valAddr,
								Status:          stakingtypes.Bonded,
								Tokens:          sdkmath.NewInt(1000000),
								DelegatorShares: sdkmath.LegacyNewDec(1000000),
							})
						}
						return nil
					})
			mocks.AcctKeeper.EXPECT().
				GetModuleAddress(shareclasstypes.ModuleName).
				Return(moduleAddr[0]).
				AnyTimes()
			mocks.AcctKeeper.EXPECT().
				GetModuleAddress(types.ModuleName).
				Return(moduleAddr[1]).
				AnyTimes()
			mocks.StakingKeeper.EXPECT().
				IterateDelegations(ctx, moduleAddr[0], gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, voter sdk.AccAddress, fn func(index int64, d stakingtypes.DelegationI) bool) error {
						return nil
					},
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

			_, err := k.StartNewEpoch(ctx, &types.MsgStartNewEpoch{
				Sender: addrs[0].String(),
			})
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
			} else {
				_, found, err := k.GetLastEpoch(ctx)
				require.NoError(t, err)
				require.False(t, found)
			}
		})
	}
}
