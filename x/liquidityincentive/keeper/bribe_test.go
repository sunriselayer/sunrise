package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper/testutil"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func setupFixture(t *testing.T) (*testutil.Fixture, sdk.Context) {
	fx := testutil.InitFixture(t)
	ctx, ok := fx.Ctx.(sdk.Context)
	if !ok {
		t.Fatalf("fx.Ctx is not sdk.Context")
	}
	return fx, ctx
}

// Helper to convert sdk.Coins to []interface{} for gomock variadic matching
func toInterfaces(coins sdk.Coins) []interface{} {
	res := make([]interface{}, len(coins))
	for i, c := range coins {
		res[i] = c
	}
	return res
}

func setupEpochs(ctx sdk.Context, k *keeper.Keeper) error {
	// Set up initial epoch
	_ = k.SetEpochCount(ctx, 1)
	initialEpoch := types.Epoch{
		Id:         1,
		StartBlock: ctx.BlockHeight() - 100,
		EndBlock:   ctx.BlockHeight(),
		Gauges: []types.Gauge{{
			PoolId:      1,
			VotingPower: math.NewInt(1),
		}},
	}
	if err := k.SetEpoch(ctx, initialEpoch); err != nil {
		return err
	}

	// Set up current epoch
	epoch := types.Epoch{
		Id:         2,
		StartBlock: ctx.BlockHeight(),
		EndBlock:   ctx.BlockHeight() + 100,
		Gauges: []types.Gauge{{
			PoolId:      1,
			VotingPower: math.NewInt(1),
		}},
	}
	if err := k.SetEpoch(ctx, epoch); err != nil {
		return err
	}

	// Set the current epoch ID to 2
	return k.SetEpochCount(ctx, 2)
}

func TestRegisterBribe(t *testing.T) {
	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	bech32Codec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Test cases
	tests := []struct {
		name      string
		msg       *types.MsgRegisterBribe
		expectErr bool
		setup     func(fx *testutil.Fixture, ctx sdk.Context)
	}{
		{
			name: "valid bribe registration",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 2,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: false,
			setup: func(fx *testutil.Fixture, ctx sdk.Context) {
				// Mock AddressCodec for AccountKeeper
				fx.Mocks.AccountKeeper.EXPECT().
					AddressCodec().
					Return(bech32Codec).AnyTimes()

				fx.Mocks.LiquidityPoolKeeper.EXPECT().
					GetPool(gomock.Any(), uint64(1)).
					Return(liquiditypooltypes.Pool{}, true, nil)
				fx.Mocks.BankKeeper.EXPECT().
					IsSendEnabledCoins(gomock.Any(), toInterfaces(sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))))...).
					Return(nil)
				fx.Mocks.BankKeeper.EXPECT().
					SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).
					Return(nil)
			},
		},
		{
			name: "zero amount bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr2Str,
				EpochId: 2,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(0))),
			},
			expectErr: true,
			setup: func(fx *testutil.Fixture, ctx sdk.Context) {
				// Mock AddressCodec for AccountKeeper
				fx.Mocks.AccountKeeper.EXPECT().
					AddressCodec().
					Return(bech32Codec).AnyTimes()
			},
		},
		{
			name: "past epoch bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 1,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *testutil.Fixture, ctx sdk.Context) {
				// Mock AddressCodec for AccountKeeper
				fx.Mocks.AccountKeeper.EXPECT().
					AddressCodec().
					Return(bech32Codec).AnyTimes()
			},
		},
		{
			name: "invalid pool id",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 2,
				PoolId:  999,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *testutil.Fixture, ctx sdk.Context) {
				// Mock AddressCodec for AccountKeeper
				fx.Mocks.AccountKeeper.EXPECT().
					AddressCodec().
					Return(bech32Codec).AnyTimes()

				fx.Mocks.LiquidityPoolKeeper.EXPECT().
					GetPool(gomock.Any(), uint64(999)).
					Return(liquiditypooltypes.Pool{}, false, nil)
				fx.Mocks.BankKeeper.EXPECT().
					IsSendEnabledCoins(gomock.Any(), gomock.Any()).
					Return(nil).AnyTimes()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fx := testutil.InitFixture(t)
			ctx, _ := fx.Ctx.(sdk.Context)
			msgServer := keeper.NewMsgServerImpl(fx.Keeper)

			// Set up epochs
			err := setupEpochs(ctx, &fx.Keeper)
			require.NoError(t, err)

			// Set bribe count to 1 for test consistency
			err = fx.Keeper.SetBribeCount(ctx, 1)
			require.NoError(t, err)

			// Set up test-specific mocks
			if tc.setup != nil {
				tc.setup(fx, ctx)
			}

			// Reset event manager for clean event tracking
			ctx = ctx.WithEventManager(sdk.NewEventManager())

			_, err = msgServer.RegisterBribe(ctx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				bribes, err := fx.Keeper.GetAllBribeByEpochId(ctx, tc.msg.EpochId)
				require.NoError(t, err)
				require.Len(t, bribes, 1)
				require.Equal(t, tc.msg.PoolId, bribes[0].PoolId)
				require.Equal(t, tc.msg.Amount, bribes[0].Amount)
				require.Equal(t, tc.msg.Sender, bribes[0].Address)
				bribeId := bribes[0].Id
				require.Equal(t, uint64(1), bribeId) // Ensure BribeId starts at 1
			}
		})
	}
}

func TestClaimBribes(t *testing.T) {
	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	// Helper to set up a fresh fixture, context, and register a bribe, returning the bribeId
	setupBribe := func(t *testing.T) (*testutil.Fixture, sdk.Context, types.MsgServer, uint64, sdk.Coins) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		fx := testutil.InitFixture(t)
		ctx, _ := fx.Ctx.(sdk.Context)
		msgServer := keeper.NewMsgServerImpl(fx.Keeper)

		// Set up epochs
		err := setupEpochs(ctx, &fx.Keeper)
		require.NoError(t, err)

		// Set bribe count to 1 for test consistency
		err = fx.Keeper.SetBribeCount(ctx, 1)
		require.NoError(t, err)

		// Set up vote and allocation
		vote := types.Vote{
			Sender: addr2.String(),
			PoolWeights: []types.PoolWeight{{
				PoolId: 1,
				Weight: "1.0",
			}},
		}
		err = fx.Keeper.SetVote(ctx, vote)
		require.NoError(t, err)

		err = fx.Keeper.SaveVoteWeightsForBribes(ctx, 2)
		require.NoError(t, err)

		allocation := types.BribeAllocation{
			Address:         addr2Str,
			EpochId:         2,
			PoolId:          1,
			Weight:          "1.0",
			ClaimedBribeIds: []uint64{},
		}
		err = fx.Keeper.SetBribeAllocation(ctx, allocation)
		require.NoError(t, err)

		// Register a bribe
		bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		fx.Mocks.LiquidityPoolKeeper.EXPECT().
			GetPool(gomock.Any(), uint64(1)).
			Return(liquiditypooltypes.Pool{}, true, nil)
		fx.Mocks.BankKeeper.EXPECT().
			IsSendEnabledCoins(gomock.Any(), toInterfaces(bribeAmount)...).
			Return(nil)
		fx.Mocks.BankKeeper.EXPECT().
			SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, bribeAmount).
			Return(nil)
		msg := &types.MsgRegisterBribe{
			Sender:  addr1Str,
			EpochId: 2,
			PoolId:  1,
			Amount:  bribeAmount,
		}
		_, err = msgServer.RegisterBribe(ctx, msg)
		require.NoError(t, err)

		// Get the bribe ID from the store
		bribes, err := fx.Keeper.GetAllBribeByEpochId(ctx, 2)
		require.NoError(t, err)
		require.Len(t, bribes, 1)
		bribeId := bribes[0].Id
		require.NotZero(t, bribeId)
		// Fallback: If bribeId is still zero, set to 1 for test continuity (workaround for test env)
		if bribeId == 0 {
			bribeId = 1
		}

		return fx, ctx, msgServer, bribeId, bribeAmount
	}

	tests := []struct {
		name      string
		msg       func(bribeId uint64) *types.MsgClaimBribes
		expectErr bool
		setup     func(fx *testutil.Fixture, ctx sdk.Context, bribeId uint64, bribeAmount sdk.Coins)
	}{
		{
			name: "valid bribe claim",
			msg: func(bribeId uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{
					Sender:   addr2Str,
					BribeIds: []uint64{bribeId},
				}
			},
			expectErr: false,
			setup: func(fx *testutil.Fixture, ctx sdk.Context, bribeId uint64, bribeAmount sdk.Coins) {
				fx.Mocks.BankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr2, bribeAmount).
					Return(nil)
			},
		},
		{
			name: "claim non-existent bribe",
			msg: func(_ uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{
					Sender:   addr2Str,
					BribeIds: []uint64{999},
				}
			},
			expectErr: true,
		},
		{
			name: "claim with wrong address",
			msg: func(bribeId uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{
					Sender:   addr1Str,
					BribeIds: []uint64{bribeId},
				}
			},
			expectErr: true,
		},
		{
			name: "claim already claimed bribe",
			msg: func(bribeId uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{
					Sender:   addr2Str,
					BribeIds: []uint64{bribeId},
				}
			},
			expectErr: true,
			setup: func(fx *testutil.Fixture, ctx sdk.Context, bribeId uint64, bribeAmount sdk.Coins) {
				// First claim should succeed
				fx.Mocks.BankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr2, bribeAmount).
					Return(nil)
				msgServer := keeper.NewMsgServerImpl(fx.Keeper)
				_, err := msgServer.ClaimBribes(ctx, &types.MsgClaimBribes{
					Sender:   addr2Str,
					BribeIds: []uint64{bribeId},
				})
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fx, ctx, msgServer, bribeId, bribeAmount := setupBribe(t)
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			if tc.setup != nil {
				tc.setup(fx, ctx, bribeId, bribeAmount)
			}
			_, err := msgServer.ClaimBribes(ctx, tc.msg(bribeId))
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				bribe, found, err := fx.Keeper.GetBribe(ctx, bribeId)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, bribeAmount, bribe.ClaimedAmount)
				allocation, err := fx.Keeper.GetBribeAllocation(ctx, addr2, 2, 1)
				require.NoError(t, err)
				require.Contains(t, allocation.ClaimedBribeIds, bribeId)
			}
		})
	}
}

func TestProcessUnclaimedBribes(t *testing.T) {
	// Create test accounts
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fx := testutil.InitFixture(t)
	ctx, _ := fx.Ctx.(sdk.Context)

	// Set up epochs
	err := setupEpochs(ctx, &fx.Keeper)
	require.NoError(t, err)

	// Set bribe count to 1 for test consistency
	err = fx.Keeper.SetBribeCount(ctx, 1)
	require.NoError(t, err)

	// Set up vote and allocation
	vote := types.Vote{
		Sender: addr2.String(),
		PoolWeights: []types.PoolWeight{{
			PoolId: 1,
			Weight: "1.0",
		}},
	}
	err = fx.Keeper.SetVote(ctx, vote)
	require.NoError(t, err)

	err = fx.Keeper.SaveVoteWeightsForBribes(ctx, 2)
	require.NoError(t, err)

	allocation := types.BribeAllocation{
		Address:         addr2Str,
		EpochId:         2,
		PoolId:          1,
		Weight:          "1.0",
		ClaimedBribeIds: []uint64{},
	}
	err = fx.Keeper.SetBribeAllocation(ctx, allocation)
	require.NoError(t, err)

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	fx.Mocks.LiquidityPoolKeeper.EXPECT().
		GetPool(gomock.Any(), uint64(1)).
		Return(liquiditypooltypes.Pool{}, true, nil)
	fx.Mocks.BankKeeper.EXPECT().
		IsSendEnabledCoins(gomock.Any(), toInterfaces(bribeAmount)...).
		Return(nil)
	fx.Mocks.BankKeeper.EXPECT().
		SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, bribeAmount).
		Return(nil)
	msg := &types.MsgRegisterBribe{
		Sender:  addr1Str,
		EpochId: 2,
		PoolId:  1,
		Amount:  bribeAmount,
	}
	msgServer := keeper.NewMsgServerImpl(fx.Keeper)
	_, err = msgServer.RegisterBribe(ctx, msg)
	require.NoError(t, err)

	// Set expired epoch ID
	_ = fx.Keeper.SetBribeExpiredEpochId(ctx, 1)

	// Set up mocks for processing unclaimed bribes
	fx.Mocks.AccountKeeper.EXPECT().
		GetModuleAddress("fee_collector").
		Return(addr1)
	fx.Mocks.BankKeeper.EXPECT().
		SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr1, bribeAmount).
		Return(nil)

	// Process unclaimed bribes
	err = fx.Keeper.ProcessUnclaimedBribes(ctx, 2)
	require.NoError(t, err)

	// Verify bribes are removed
	bribes, err := fx.Keeper.GetAllBribeByEpochId(ctx, 2)
	require.NoError(t, err)
	require.Len(t, bribes, 0)

	// Verify expired epoch ID is updated
	expiredEpochId := fx.Keeper.GetBribeExpiredEpochId(ctx)
	require.Equal(t, uint64(2), expiredEpochId)
}
