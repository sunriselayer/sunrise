package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
	liquiditypooltypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func setupEpochsForBribe(ctx sdk.Context, k *keeper.Keeper) error {
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

	// Set up next epoch
	nextEpoch := types.Epoch{
		Id:         3,
		StartBlock: ctx.BlockHeight() + 101,
		EndBlock:   ctx.BlockHeight() + 200,
		Gauges: []types.Gauge{{
			PoolId:      1,
			VotingPower: math.NewInt(1),
		}},
	}
	if err := k.SetEpoch(ctx, nextEpoch); err != nil {
		return err
	}

	// Set the current epoch ID to 3
	return k.SetEpochCount(ctx, 3)
}

func TestProcessUnclaimedBribes(t *testing.T) {
	_, _, addr1 := testdata.KeyTestPubAddr()
	addr1Str := addr1.String()
	_, _, addr2 := testdata.KeyTestPubAddr()
	addr2Str := addr2.String()

	bech32Codec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fx := initFixture(t)
	sdkCtx := fx.ctx.(sdk.Context)
	msgServer := keeper.NewMsgServerImpl(fx.keeper)

	// Set up mocks for bribe registration first
	fx.mocks.AcctKeeper.EXPECT().AddressCodec().Return(bech32Codec).AnyTimes()
	fx.mocks.LiquiditypoolKeeper.EXPECT().GetPool(gomock.Any(), uint64(1)).Return(liquiditypooltypes.Pool{}, true, nil).AnyTimes()
	fx.mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()
	fx.mocks.BankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.BribeAccount, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))).Return(nil).AnyTimes()

	// Now set up epochs and bribe count
	sdkCtx = sdkCtx.WithBlockHeight(1000)
	err := setupEpochsForBribe(sdkCtx, &fx.keeper)
	require.NoError(t, err)
	err = fx.keeper.SetBribeCount(sdkCtx, 1)
	require.NoError(t, err)

	// Set up vote and allocation
	vote := types.Vote{Sender: addr2Str, PoolWeights: []types.PoolWeight{{PoolId: 1, Weight: "1.0"}}}
	err = fx.keeper.SetVote(sdkCtx, vote)
	require.NoError(t, err)
	err = fx.keeper.SaveVoteWeightsForBribes(sdkCtx, 4)
	require.NoError(t, err)
	allocation := types.BribeAllocation{Address: addr2Str, EpochId: 4, PoolId: 1, Weight: "1.0", ClaimedBribeIds: []uint64{}}
	err = fx.keeper.SetBribeAllocation(sdkCtx, allocation)
	require.NoError(t, err)

	// Register a bribe
	bribeAmount := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))
	sdkCtx = sdkCtx.WithEventManager(sdk.NewEventManager())
	msg := &types.MsgRegisterBribe{Sender: addr1Str, EpochId: 4, PoolId: 1, Amount: bribeAmount}
	_, err = msgServer.RegisterBribe(sdkCtx, msg)
	require.NoError(t, err)

	// Set expired epoch ID
	_ = fx.keeper.SetBribeExpiredEpochId(sdkCtx, 2)

	// Set up mocks for processing unclaimed bribes
	fx.mocks.AcctKeeper.EXPECT().GetModuleAddress("fee_collector").Return(addr1)
	fx.mocks.BankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.BribeAccount, addr1, bribeAmount).Return(nil)

	// Process unclaimed bribes
	err = fx.keeper.ProcessUnclaimedBribes(sdkCtx, 4)
	require.NoError(t, err)

	// Verify bribes are removed
	bribes, err := fx.keeper.GetAllBribeByEpochId(sdkCtx, 4)
	require.NoError(t, err)
	require.Len(t, bribes, 0)

	// Verify expired epoch ID is updated
	expiredEpochId := fx.keeper.GetBribeExpiredEpochId(sdkCtx)
	require.Equal(t, uint64(4), expiredEpochId)
}
