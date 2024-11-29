package keeper_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"

	"github.com/cosmos/cosmos-sdk/codec/address"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type tallyFixture struct {
	t        *testing.T
	valAddrs []sdk.ValAddress
	delAddrs []sdk.AccAddress
	keeper   *keeper.Keeper
	ctx      sdk.Context
	mocks    keepertest.LiquidityIncentiveMocks
}

var (
	// handy functions
	setTotalBonded = func(s tallyFixture, n int64) {
		s.mocks.AcctKeeper.EXPECT().AddressCodec().Return(address.NewBech32Codec("sunrise")).AnyTimes()
		s.mocks.StakingKeeper.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec("sunrisevaloper")).AnyTimes()
		s.mocks.StakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(sdkmath.NewInt(n), nil)
	}
	delegatorVote = func(s tallyFixture, voter sdk.AccAddress, delegations []stakingtypes.Delegation, weights []types.PoolWeight) {
		s.keeper.SetVote(s.ctx, types.Vote{
			Sender:      voter.String(),
			PoolWeights: weights,
		})
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
		expectedTally []types.TallyResult
		expectError   bool
	}{
		{
			name: "no votes, no bonded tokens: tally fails",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 0)
			},
			expectedTally: []types.TallyResult{},
		},
		{
			name: "no votes: tally fails",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
			},
			expectedTally: []types.TallyResult{},
		},
		{
			name: "one validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(1000000)}},
		},
		{
			name: "one account votes without delegation",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				delegatorVote(s, s.delAddrs[0], nil, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.TallyResult{},
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
					Shares:           sdkmath.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(42)}},
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
					Shares:           sdkmath.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(1000000)}},
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
						Shares:           sdkmath.LegacyNewDec(21),
					},
					{
						DelegatorAddress: del0Addr,
						ValidatorAddress: val1Addr,
						Shares:           sdkmath.LegacyNewDec(21),
					},
				}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[1], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
				validatorVote(s, s.valAddrs[2], []types.PoolWeight{{PoolId: 1, Weight: "1"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(3000000)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, mocks, ctx := keepertest.LiquidityincentiveKeeper(t)

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
					func(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) bool) error {
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
				ctx:      ctx,
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			tallyWeights, err := k.Tally(ctx)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedTally, tallyWeights)
		})
	}
}

func TestTally_MultipleChoice(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(tallyFixture)
		expectedTally []types.TallyResult
		expectError   bool
	}{
		{
			name: "no votes, no bonded tokens",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 0)
			},
			expectedTally: []types.TallyResult{},
		},
		{
			name: "no votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
			},
			expectedTally: []types.TallyResult{},
		},
		{
			name: "one validator votes",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(500000)}, {PoolId: 2, Count: math.NewIntFromUint64(500000)}},
		},
		{
			name: "one account votes without delegation",
			setup: func(s tallyFixture) {
				setTotalBonded(s, 10000000)
				delegatorVote(s, s.delAddrs[0], nil, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.TallyResult{},
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
					Shares:           sdkmath.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(21)}, {PoolId: 2, Count: math.NewIntFromUint64(21)}},
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
					Shares:           sdkmath.LegacyNewDec(42),
				}}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(500000)}, {PoolId: 2, Count: math.NewIntFromUint64(500000)}},
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
						Shares:           sdkmath.LegacyNewDec(21),
					},
					{
						DelegatorAddress: del0Addr,
						ValidatorAddress: val1Addr,
						Shares:           sdkmath.LegacyNewDec(21),
					},
				}
				delegatorVote(s, s.delAddrs[0], delegations, []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[0], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[1], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
				validatorVote(s, s.valAddrs[2], []types.PoolWeight{{PoolId: 1, Weight: "0.5"}, {PoolId: 2, Weight: "0.5"}})
			},
			expectedTally: []types.TallyResult{{PoolId: 1, Count: math.NewIntFromUint64(1500000)}, {PoolId: 2, Count: math.NewIntFromUint64(1500000)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, mocks, ctx := keepertest.LiquidityincentiveKeeper(t)
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
					func(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) bool) error {
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
				ctx:      ctx,
				keeper:   &k,
				mocks:    mocks,
			}
			tt.setup(suite)

			tallyWeights, err := k.Tally(ctx)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedTally, tallyWeights)
		})
	}
}
