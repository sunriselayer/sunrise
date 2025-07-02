package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sunriselayer/sunrise/app/consts"
	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestTickInfoStore(t *testing.T) {
	f := initFixture(t)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	k := f.keeper

	// Not available tick
	_, err := k.GetTickInfo(ctx, 1, 1)
	require.Error(t, err)

	feeGrowth := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	k.SetTickInfo(ctx, types.TickInfo{
		PoolId:         1,
		TickIndex:      1,
		LiquidityGross: math.LegacyOneDec().String(),
		LiquidityNet:   math.LegacyOneDec().String(),
		FeeGrowth:      feeGrowth,
	})
	k.SetTickInfo(ctx, types.TickInfo{
		PoolId:         1,
		TickIndex:      2,
		LiquidityGross: math.LegacyOneDec().String(),
		LiquidityNet:   math.LegacyOneDec().String(),
		FeeGrowth:      feeGrowth,
	})
	k.SetTickInfo(ctx, types.TickInfo{
		PoolId:         2,
		TickIndex:      1,
		LiquidityGross: math.LegacyOneDec().String(),
		LiquidityNet:   math.LegacyOneDec().String(),
		FeeGrowth:      feeGrowth,
	})

	tickInfo, err := k.GetTickInfo(ctx, 1, 1)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(1))
	require.Equal(t, tickInfo.FeeGrowth.String(), "1.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross, "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "1.000000000000000000")

	tickInfo, err = k.GetTickInfo(ctx, 2, 1)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(2))
	require.Equal(t, tickInfo.TickIndex, int64(1))
	require.Equal(t, tickInfo.FeeGrowth.String(), "1.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross, "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "1.000000000000000000")

	tickInfo, err = k.GetTickInfo(ctx, 1, 2)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(2))
	require.Equal(t, tickInfo.FeeGrowth.String(), "1.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross, "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "1.000000000000000000")

	tickInfos := k.GetAllInitializedTicksForPool(ctx, 1)
	require.Len(t, tickInfos, 2)

	tickInfos = k.GetAllInitializedTicksForPool(ctx, 2)
	require.Len(t, tickInfos, 1)

	tickInfos = k.GetAllTickInfos(ctx)
	require.Len(t, tickInfos, 3)

	k.RemoveTickInfo(ctx, 1, 2)
	_, err = k.GetTickInfo(ctx, 1, 2)
	require.Error(t, err)

	tickInfos = k.GetAllTickInfos(ctx)
	require.Len(t, tickInfos, 2)
}

func TestUpsertTick(t *testing.T) {
	k, mocks, srv, ctx := setupMsgServer(t)
	quoteDenom := consts.StableDenom

	// When pool does not exist
	_, err := k.UpsertTick(ctx, 1, 0, math.LegacyNewDec(10), true)
	require.Error(t, err)

	// When pool exist
	wctx := sdk.UnwrapSDKContext(ctx)

	mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	sender := sdk.AccAddress("sender")
	_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
		Sender:     sender.String(),
		DenomBase:  "base",
		DenomQuote: quoteDenom,
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "-0.5",
	})
	require.NoError(t, err)

	tickEmpty, err := k.UpsertTick(ctx, 0, 0, math.LegacyNewDec(10), true)
	require.NoError(t, err)
	require.False(t, tickEmpty)

	// Check state change in tickInfo
	tickInfo, err := k.GetTickInfo(ctx, 0, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(0))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross, "10.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "-10.000000000000000000")

	// Tick's available
	tickEmpty, err = k.UpsertTick(ctx, 0, 0, math.LegacyNewDec(10), false)
	require.NoError(t, err)
	require.False(t, tickEmpty)

	// Check state change in tickInfo
	tickInfo, err = k.GetTickInfo(ctx, 0, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(0))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross, "20.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "0.000000000000000000")

	// Negative deltaLiquidity
	tickEmpty, err = k.UpsertTick(ctx, 0, 0, math.LegacyNewDec(-20), false)
	require.NoError(t, err)
	require.False(t, tickEmpty)

	// Check state change in tickInfo
	tickInfo, err = k.GetTickInfo(ctx, 0, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(0))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross, "0.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "-20.000000000000000000")
}

func TestNewTickInfo(t *testing.T) {
	k, mocks, srv, ctx := setupMsgServer(t)
	quoteDenom := consts.StableDenom
	// When pool does not exist
	_, err := k.NewTickInfo(ctx, 1, 0)
	require.Error(t, err)

	// When empty pool exist
	wctx := sdk.UnwrapSDKContext(ctx)

	mocks.BankKeeper.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mocks.BankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	sender := sdk.AccAddress("sender")
	_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
		Sender:     sender.String(),
		DenomBase:  "base",
		DenomQuote: quoteDenom,
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "-0.5",
	})
	require.NoError(t, err)

	tickInfo, err := k.NewTickInfo(ctx, 0, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(0))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross, "0.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "0.000000000000000000")

	// When pool accumulator has positive accumulation value
	accumulator, err := k.GetFeeAccumulator(ctx, 0)
	require.NoError(t, err)
	accumulator.AccumValue = sdk.DecCoins{sdk.NewInt64DecCoin("denom", 100)}
	err = k.SetAccumulator(ctx, accumulator)
	require.NoError(t, err)

	tickInfo, err = k.NewTickInfo(ctx, 0, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(0))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "100.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross, "0.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "0.000000000000000000")
}

func TestDecodeTickBytes(t *testing.T) {
	// When tick bytes' empty
	_, err := keeper.DecodeTickBytes([]byte{})
	require.Error(t, err)

	// When tick bytes' invalid
	_, err = keeper.DecodeTickBytes([]byte{0x1})
	require.Error(t, err)

	// Valid tick bytes
	tickInfo := types.TickInfo{
		PoolId:         1,
		TickIndex:      0,
		LiquidityGross: math.LegacyOneDec().String(),
		LiquidityNet:   math.LegacyOneDec().String(),
		FeeGrowth:      sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1))),
	}
	bz, err := tickInfo.Marshal()
	require.NoError(t, err)

	decoded, err := keeper.DecodeTickBytes(bz)
	require.NoError(t, err)
	require.Equal(t, tickInfo, decoded)
}

func TestCrossTick(t *testing.T) {
	f := initFixture(t)
	ctx := sdk.UnwrapSDKContext(f.ctx)
	k := f.keeper

	oneDecCoins := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	twoDecCoins := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(2)))
	threeDecCoins := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(3)))

	// When tickInfo's empty
	err := k.CrossTick(ctx, 1, 0, nil, oneDecCoins[0], twoDecCoins)
	require.Error(t, err)

	// When tickInfo's valid
	err = k.CrossTick(ctx, 1, 0, &types.TickInfo{
		PoolId:         1,
		TickIndex:      0,
		LiquidityGross: math.LegacyOneDec().String(),
		LiquidityNet:   math.LegacyOneDec().String(),
		FeeGrowth:      twoDecCoins,
	}, oneDecCoins[0], threeDecCoins)
	require.NoError(t, err)

	// check TickInfo update
	tickInfo, err := k.GetTickInfo(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "2.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross, "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet, "1.000000000000000000")
}
