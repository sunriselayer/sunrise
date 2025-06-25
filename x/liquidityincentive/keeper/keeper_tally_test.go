package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"cosmossdk.io/math"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type tallyFixture struct {
	t        *testing.T
	valAddrs []sdk.ValAddress
	delAddrs []sdk.AccAddress
	keeper   *keeper.Keeper
	ctx      sdk.Context
	mocks    LiquidityIncentiveMocks
}

var (
	// handy functions
	setTotalBonded = func(s tallyFixture, n int64) {
		s.mocks.AcctKeeper.EXPECT().
			AddressCodec().
			Return(addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())).AnyTimes()
		s.mocks.StakingKeeper.EXPECT().
			ValidatorAddressCodec().
			Return(addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())).AnyTimes()
		s.mocks.StakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(math.NewInt(n), nil)
	}
	delegatorVote = func(s tallyFixture, voter sdk.AccAddress, delegations []stakingtypes.Delegation, weights []types.PoolWeight) {
		err := s.keeper.SetVote(s.ctx, types.Vote{
			Sender:      voter.String(),
			PoolWeights: weights,
		})
		require.NoError(s.t, err)
		s.mocks.StakingKeeper.EXPECT().
			IterateDelegations(s.ctx, voter, gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, voter sdk.AccAddress, fn func(index int64, d stakingtypes.DelegationI) bool) error {
					for i, d := range delegations {
						fn(int64(i), d)
					}
					return nil
				})
	}
	validatorVote = func(s tallyFixture, voter sdk.ValAddress, weights []types.PoolWeight) {
		// validatorVote is like delegatorVote but without delegations
		delegatorVote(s, sdk.AccAddress(voter), nil, weights)
	}
)

func TestTally_Standard(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(tallyFixture)
		expectedTally []types.Gauge
		expectError   bool
	}{
		{
			name: "no votes, no bonded tokens: tally fails",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 0)
			},
			expectedTally: []types.Gauge{},
		},
		{
			name: "no votes: tally fails",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
			},
			expectedTally: []types.Gauge{},
		},
		{
			name: "one validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(1000000)}},
		},
		{
			name: "one account votes without delegation",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				delegatorVote(s, s.delAddrs[0], nil, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.Gauge{},
		},
		{
			name: "one delegator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				del0Addr, err := s.mocks.AcctKeeper.AddressCodec().BytesToString(s.delAddrs[0])
				require.NoError(t, err)
				val0Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[0])
				require.NoError(t, err)
				delegations := []stakingtypes.Delegation{{
					DelegatorAddress: del0Addr,
					ValidatorAddress: val0Addr,
					Shares:           math.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(42)}},
		},
		{
			name: "one delegator votes, validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				del0Addr, err := s.mocks.AcctKeeper.AddressCodec().BytesToString(s.delAddrs[0])
				require.NoError(t, err)
				val0Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[0])
				require.NoError(t, err)
				delegations := []stakingtypes.Delegation{{
					DelegatorAddress: del0Addr,
					ValidatorAddress: val0Addr,
					Shares:           math.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(1000000)}},
		},
		{
			name: "delegator with mixed delegations",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				del0Addr, err := s.mocks.AcctKeeper.AddressCodec().BytesToString(s.delAddrs[0])
				require.NoError(t, err)
				val0Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[0])
				require.NoError(t, err)
				val1Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[1])
				require.NoError(t, err)
				delegations := []stakingtypes.Delegation{
					{
						DelegatorAddress: del0Addr,
						ValidatorAddress: val0Addr,
						Shares:           math.LegacyNewDec(21),
					},
					{
						DelegatorAddress: del0Addr,
						ValidatorAddress: val1Addr,
						Shares:           math.LegacyNewDec(21),
					},
				}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[1], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[2], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(3000000)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := sdk.UnwrapSDKContext(f.ctx)
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
							valAddr, err := mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(valAddrs[i])
							require.NoError(t, err)
							fn(i, stakingtypes.Validator{
								OperatorAddress: valAddr,
								Status:          stakingtypes.Bonded,
								Tokens:          math.NewInt(1000000),
								DelegatorShares: math.LegacyNewDec(1000000),
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
				ctx:      ctx,
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			_, gauges, err := k.Tally(ctx)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedTally, gauges)
		})
	}
}

func TestTally_MultipleChoice(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(tallyFixture)
		expectedTally []types.Gauge
		expectError   bool
	}{
		{
			name: "no votes, no bonded tokens",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 0)
			},
			expectedTally: []types.Gauge{},
		},
		{
			name: "no votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
			},
			expectedTally: []types.Gauge{},
		},
		{
			name: "one validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(500000)}, {PoolId: 2, VotingPower: math.NewIntFromUint64(500000)}},
		},
		{
			name: "one account votes without delegation",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				delegatorVote(s, s.delAddrs[0], nil, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.Gauge{},
		},
		{
			name: "one delegator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				del0Addr, err := s.mocks.AcctKeeper.AddressCodec().BytesToString(s.delAddrs[0])
				require.NoError(t, err)
				val0Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[0])
				require.NoError(t, err)
				delegations := []stakingtypes.Delegation{{
					DelegatorAddress: del0Addr,
					ValidatorAddress: val0Addr,
					Shares:           math.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(21)}, {PoolId: 2, VotingPower: math.NewIntFromUint64(21)}},
		},
		{
			name: "one delegator votes, validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				del0Addr, err := s.mocks.AcctKeeper.AddressCodec().BytesToString(s.delAddrs[0])
				require.NoError(t, err)
				val0Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[0])
				require.NoError(t, err)
				delegations := []stakingtypes.Delegation{{
					DelegatorAddress: del0Addr,
					ValidatorAddress: val0Addr,
					Shares:           math.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(500000)}, {PoolId: 2, VotingPower: math.NewIntFromUint64(500000)}},
		},
		{
			name: "delegator with mixed delegations",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				del0Addr, err := s.mocks.AcctKeeper.AddressCodec().BytesToString(s.delAddrs[0])
				require.NoError(t, err)
				val0Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[0])
				require.NoError(t, err)
				val1Addr, err := s.mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(s.valAddrs[1])
				require.NoError(t, err)
				delegations := []stakingtypes.Delegation{
					{
						DelegatorAddress: del0Addr,
						ValidatorAddress: val0Addr,
						Shares:           math.LegacyNewDec(21),
					},
					{
						DelegatorAddress: del0Addr,
						ValidatorAddress: val1Addr,
						Shares:           math.LegacyNewDec(21),
					},
				}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[1], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[2], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.Gauge{{PoolId: 1, VotingPower: math.NewIntFromUint64(1500000)}, {PoolId: 2, VotingPower: math.NewIntFromUint64(1500000)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := sdk.UnwrapSDKContext(f.ctx)
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
							valAddr, err := mocks.StakingKeeper.ValidatorAddressCodec().BytesToString(valAddrs[i])
							require.NoError(t, err)
							fn(i, stakingtypes.Validator{
								OperatorAddress: valAddr,
								Status:          stakingtypes.Bonded,
								Tokens:          math.NewInt(1000000),
								DelegatorShares: math.LegacyNewDec(1000000),
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
				ctx:      ctx,
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			_, gauges, err := k.Tally(ctx)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedTally, gauges)
		})
	}
}
