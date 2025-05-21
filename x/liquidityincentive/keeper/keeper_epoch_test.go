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
			expectError:   true,
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
			f := initFixture(t)
			ctx := f.ctx
			k := f.keeper
			mocks := f.mocks

			var (
				numVals       = 10
				numDelegators = 5
				moduleAddr    = simtestutil.CreateRandomAccounts(2)
				addrs         = simtestutil.CreateRandomAccounts(numVals + numDelegators)
				valAddrs      = simtestutil.ConvertAddrsToValAddrs(addrs[:numVals])
				delAddrs      = addrs[numVals:]
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

			err := k.CreateEpoch(sdk.UnwrapSDKContext(ctx), 1)
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
			}
		})
	}
}
