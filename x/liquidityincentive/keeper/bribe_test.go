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
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

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
	currentEpoch := types.Epoch{
		Id:         2,
		StartBlock: ctx.BlockHeight(),
		EndBlock:   ctx.BlockHeight() + 100,
		Gauges: []types.Gauge{{
			PoolId:      1,
			VotingPower: math.NewInt(1),
		}},
	}
	if err := k.SetEpoch(ctx, currentEpoch); err != nil {
		return err
	}

	// Set up future epoch for bribes
	futureEpoch := types.Epoch{
		Id:         3,
		StartBlock: ctx.BlockHeight() + 101,
		EndBlock:   ctx.BlockHeight() + 200,
		Gauges: []types.Gauge{{
			PoolId:      1,
			VotingPower: math.NewInt(1),
		}},
	}
	if err := k.SetEpoch(ctx, futureEpoch); err != nil {
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
		setup     func(fx *fixture, ctx sdk.Context)
	}{
		{
			name: "valid bribe registration",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 3, // Future epoch
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: false,
			setup: func(fx *fixture, ctx sdk.Context) {
				// Set up mocks first
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(1)).Return(liquiditypooltypes.Pool{}, true, nil)
				fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil)
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil)
				// Now set up epochs and bribe count
				err := setupEpochs(ctx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(ctx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "zero amount bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr2Str,
				EpochId: 3,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(0))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				err := setupEpochs(ctx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(ctx, 1)
				require.NoError(t, err)
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
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				err := setupEpochs(ctx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(ctx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "current epoch bribe",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 2,
				PoolId:  1,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				err := setupEpochs(ctx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(ctx, 1)
				require.NoError(t, err)
			},
		},
		{
			name: "invalid pool id",
			msg: &types.MsgRegisterBribe{
				Sender:  addr1Str,
				EpochId: 3,
				PoolId:  999,
				Amount:  sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context) {
				fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
				fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(999)).Return(liquiditypooltypes.Pool{}, false, nil)
				fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
				err := setupEpochs(ctx, &fx.keeper)
				require.NoError(t, err)
				err = fx.keeper.SetBribeCount(ctx, 1)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fx := initFixture(t)
			sdkCtx := fx.ctx.(sdk.Context)
			msgServer := keeper.NewMsgServerImpl(fx.keeper)

			// Set up test-specific mocks and keeper state
			if tc.setup != nil {
				tc.setup(fx, sdkCtx)
			}

			// Reset event manager for clean event tracking
			sdkCtx = sdkCtx.WithEventManager(sdk.NewEventManager())

			_, err := msgServer.RegisterBribe(sdkCtx, tc.msg)
			if tc.expectErr {
				require.Error(t, err)
				if err.Error() == "amount cannot be zero" {
					require.ErrorIs(t, err, types.ErrInvalidBribe)
				} else {
					require.ErrorIs(t, err, types.ErrBribeCannotBeCreated)
					require.Contains(t, err.Error(), "Cannot be created in the current epoch. The process has already been completed.")
				}
			} else {
				require.NoError(t, err)
				bribes, err := fx.keeper.GetAllBribeByEpochId(sdkCtx, tc.msg.EpochId)
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

	bech32Codec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Helper to set up a fresh fixture, context, and register a bribe, returning the bribeId
	setupBribe := func(fx *fixture, sdkCtx sdk.Context, msgServer types.MsgServer, addr1Str, addr2Str string, bech32Codec *addresscodec.Bech32Codec) (uint64, sdk.Coins) {
		// Set up mocks for bribe registration first
		fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
		fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(1)).Return(liquiditypooltypes.Pool{}, true, nil)
		fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil)
		fx.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil)

		// Now set up epochs and bribe count
		err := setupEpochs(sdkCtx, &fx.keeper)
		require.NoError(t, err)
		err = fx.keeper.SetBribeCount(sdkCtx, 1)
		require.NoError(t, err)

		// Set up vote and allocation
		vote := types.Vote{Sender: addr2Str, PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "1.0"}}}
		err = fx.keeper.SetVote(sdkCtx, vote)
		require.NoError(t, err)
		err = fx.keeper.SaveVoteWeightsForBribes(sdkCtx, 3)
		require.NoError(t, err)
		allocation := types.BribeAllocation{Address: addr2Str, EpochId: 3, PoolId: 1, Weight: "1.0", ClaimedBribeIds: []uint64{}}
		err = fx.keeper.SetBribeAllocation(sdkCtx, allocation)
		require.NoError(t, err)

		// Register a bribe
		bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
		sdkCtx = sdkCtx.WithEventManager(sdk.NewEventManager())
		msg := &types.MsgRegisterBribe{Sender: addr1Str, EpochId: 3, PoolId: 1, Amount: bribeAmount}
		_, err = msgServer.RegisterBribe(sdkCtx, msg)
		require.NoError(t, err)
		bribes, err := fx.keeper.GetAllBribeByEpochId(sdkCtx, 3)
		require.NoError(t, err)
		require.Len(t, bribes, 1)
		bribeId := bribes[0].Id
		require.NotZero(t, bribeId)
		return bribeId, bribeAmount
	}

	tests := []struct {
		name      string
		msg       func(bribeId uint64) *types.MsgClaimBribes
		expectErr bool
		setup     func(fx *fixture, ctx sdk.Context, bribeId uint64, bribeAmount sdk.Coins)
	}{
		{
			name: "valid bribe claim",
			msg: func(bribeId uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{Sender: addr2Str, BribeIds: []uint64{bribeId}}
			},
			expectErr: false,
			setup: func(fx *fixture, ctx sdk.Context, bribeId uint64, bribeAmount sdk.Coins) {
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr2, bribeAmount).Return(nil)
			},
		},
		{
			name: "claim non-existent bribe",
			msg: func(_ uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{Sender: addr2Str, BribeIds: []uint64{999}}
			},
			expectErr: true,
		},
		{
			name: "claim with wrong address",
			msg: func(bribeId uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{Sender: addr1Str, BribeIds: []uint64{bribeId}}
			},
			expectErr: true,
		},
		{
			name: "claim already claimed bribe",
			msg: func(bribeId uint64) *types.MsgClaimBribes {
				return &types.MsgClaimBribes{Sender: addr2Str, BribeIds: []uint64{bribeId}}
			},
			expectErr: true,
			setup: func(fx *fixture, ctx sdk.Context, bribeId uint64, bribeAmount sdk.Coins) {
				fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr2, bribeAmount).Return(nil)
				msgServer := keeper.NewMsgServerImpl(fx.keeper)
				_, err := msgServer.ClaimBribes(ctx, &types.MsgClaimBribes{Sender: addr2Str, BribeIds: []uint64{bribeId}})
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fx := initFixture(t)
			sdkCtx := fx.ctx.(sdk.Context)
			msgServer := keeper.NewMsgServerImpl(fx.keeper)

			bribeId, bribeAmount := setupBribe(fx, sdkCtx, msgServer, addr1Str, addr2Str, bech32Codec.(*addresscodec.Bech32Codec))
			sdkCtx = sdkCtx.WithEventManager(sdk.NewEventManager())
			if tc.setup != nil {
				tc.setup(fx, sdkCtx, bribeId, bribeAmount)
			}
			_, err := msgServer.ClaimBribes(sdkCtx, tc.msg(bribeId))
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				bribe, found, err := fx.keeper.GetBribe(sdkCtx, bribeId)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, bribeAmount, bribe.ClaimedAmount)
				allocation, err := fx.keeper.GetBribeAllocation(sdkCtx, addr2, 3, 1)
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

	bech32Codec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fx := initFixture(t)
	sdkCtx := fx.ctx.(sdk.Context)

	// Set up mocks for bribe registration first
	fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
	fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(1)).Return(liquiditypooltypes.Pool{}, true, nil)
	fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil)
	fx.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil)

	// Now set up epochs and bribe count
	err := setupEpochs(sdkCtx, &fx.keeper)
	require.NoError(t, err)
	err = fx.keeper.SetBribeCount(sdkCtx, 1)
	require.NoError(t, err)

	// Set up vote and allocation
	vote := types.Vote{Sender: addr2.String(), PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "1.0"}}}
	err = fx.keeper.SetVote(sdkCtx, vote)
	require.NoError(t, err)
	err = fx.keeper.SaveVoteWeightsForBribes(sdkCtx, 3)
	require.NoError(t, err)
	allocation := types.BribeAllocation{Address: addr2Str, EpochId: 3, PoolId: 1, Weight: "1.0", ClaimedBribeIds: []uint64{}}
	err = fx.keeper.SetBribeAllocation(sdkCtx, allocation)
	require.NoError(t, err)

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	sdkCtx = sdkCtx.WithEventManager(sdk.NewEventManager())
	msg := &types.MsgRegisterBribe{Sender: addr1Str, EpochId: 3, PoolId: 1, Amount: bribeAmount}
	msgServer := keeper.NewMsgServerImpl(fx.keeper)
	_, err = msgServer.RegisterBribe(sdkCtx, msg)
	require.NoError(t, err)

	// Set expired epoch ID
	_ = fx.keeper.SetBribeExpiredEpochId(sdkCtx, 2)

	// Set up mocks for processing unclaimed bribes
	fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(addr1)
	fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr1, bribeAmount).Return(nil)

	// Process unclaimed bribes
	err = fx.keeper.ProcessUnclaimedBribes(sdkCtx, 3)
	require.NoError(t, err)

	// Verify bribes are removed
	bribes, err := fx.keeper.GetAllBribeByEpochId(sdkCtx, 3)
	require.NoError(t, err)
	require.Len(t, bribes, 0)

	// Verify expired epoch ID is updated
	expiredEpochId := fx.keeper.GetBribeExpiredEpochId(sdkCtx)
	require.Equal(t, uint64(3), expiredEpochId)
}
