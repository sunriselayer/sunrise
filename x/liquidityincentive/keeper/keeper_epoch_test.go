package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	shareclasstypes "github.com/sunriselayer/sunrise/x/shareclass/types"
)

func TestCreateEpoch(t *testing.T) {
	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, moduleAddr := testdata.KeyTestPubAddr()

	bech32Codec := address.NewBech32Codec("cosmos")

	tests := []struct {
		name          string
		setup         func(fx *fixture, ctx sdk.Context)
		expectedTally []types.Gauge
		expectError   bool
	}{
		{
			name: "no votes",
			setup: func(fx *fixture, ctx sdk.Context) {
				// Mock GetModuleAddress for shareclass module
				fx.mocks.AcctKeeper.EXPECT().
					GetModuleAddress(shareclasstypes.ModuleName).
					Return(moduleAddr).
					AnyTimes()

				// Mock IterateBondedValidatorsByPower to return no validators
				fx.mocks.StakingKeeper.EXPECT().
					IterateBondedValidatorsByPower(gomock.Any(), gomock.Any()).
					Return(nil)

				// Mock IterateDelegations and ValidatorAddressCodec
				fx.mocks.StakingKeeper.EXPECT().
					IterateDelegations(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).AnyTimes()
				fx.mocks.StakingKeeper.EXPECT().
					ValidatorAddressCodec().
					Return(bech32Codec).AnyTimes()
				fx.mocks.StakingKeeper.EXPECT().
					TotalBondedTokens(gomock.Any()).
					Return(math.NewInt(0), nil).AnyTimes()

				// Mock AddressCodec for AccountKeeper
				fx.mocks.AcctKeeper.EXPECT().
					AddressCodec().
					Return(bech32Codec).AnyTimes()

				// Expect FinalizeBribeForEpoch to be called, which might call other things.
				// For simplicity, we mock the dependencies of the functions called within FinalizeBribeForEpoch.
				// This test focuses on CreateEpoch's logic, assuming FinalizeBribeForEpoch is tested elsewhere.
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(sdk.AccAddress("fee_collector")).AnyTimes()
			},
			expectedTally: []types.Gauge{},
			expectError:   true,
		},
		{
			name: "one validator votes",
			setup: func(fx *fixture, ctx sdk.Context) {
				// Mock GetModuleAddress for shareclass module
				fx.mocks.AcctKeeper.EXPECT().
					GetModuleAddress(shareclasstypes.ModuleName).
					Return(moduleAddr).
					AnyTimes()

				// Mock IterateBondedValidatorsByPower to return one validator
				fx.mocks.StakingKeeper.EXPECT().
					IterateBondedValidatorsByPower(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) bool) error {
						fn(0, stakingtypes.Validator{
							OperatorAddress: addr1Str,
							Status:          stakingtypes.Bonded,
							Tokens:          math.NewInt(1000000),
							DelegatorShares: math.LegacyNewDec(1000000),
						})
						return nil
					})

				// Mock IterateDelegations and ValidatorAddressCodec
				fx.mocks.StakingKeeper.EXPECT().
					IterateDelegations(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).AnyTimes()
				fx.mocks.StakingKeeper.EXPECT().
					ValidatorAddressCodec().
					Return(bech32Codec).AnyTimes()
				fx.mocks.StakingKeeper.EXPECT().
					TotalBondedTokens(gomock.Any()).
					Return(math.NewInt(1000000), nil).AnyTimes()

				// Mock AddressCodec for AccountKeeper
				fx.mocks.AcctKeeper.EXPECT().
					AddressCodec().
					Return(bech32Codec).AnyTimes()
				
				// Expect FinalizeBribeForEpoch to be called
				fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(sdk.AccAddress("fee_collector")).AnyTimes()


				// Set up a vote with a valid address
				vote := types.Vote{
					Sender: addr1Str,
					PoolWeights: []types.PoolWeight{{
						PoolId: 1,
						Weight: "1.0",
					}},
				}
				err := fx.keeper.SetVote(ctx, vote)
				require.NoError(t, err)
			},
			expectedTally: []types.Gauge{{
				PoolId:      1,
				VotingPower: math.NewInt(1000000),
			}},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fx := initFixture(t)
			ctx, _ := fx.ctx.(sdk.Context)

			if tc.setup != nil {
				tc.setup(fx, ctx)
			}

			err := fx.keeper.CreateEpoch(ctx, 1)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				epoch, found, err := fx.keeper.GetLastEpoch(ctx)
				require.NoError(t, err)
				require.True(t, found)
				require.NotNil(t, epoch)
				require.Equal(t, tc.expectedTally, epoch.Gauges)
			}
		})
	}
}
