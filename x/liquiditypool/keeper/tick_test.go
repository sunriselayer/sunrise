package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestTickInfoStore(t *testing.T) {
	k, _, ctx := keepertest.LiquiditypoolKeeper(t)

	// Not available tick
	_, err := k.GetTickInfo(ctx, 1, 1)
	require.Error(t, err)

	feeGrowth := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	k.SetTickInfo(ctx, types.TickInfo{
		PoolId:         1,
		TickIndex:      1,
		LiquidityGross: math.LegacyOneDec(),
		LiquidityNet:   math.LegacyOneDec(),
		FeeGrowth:      feeGrowth,
	})
	k.SetTickInfo(ctx, types.TickInfo{
		PoolId:         1,
		TickIndex:      2,
		LiquidityGross: math.LegacyOneDec(),
		LiquidityNet:   math.LegacyOneDec(),
		FeeGrowth:      feeGrowth,
	})
	k.SetTickInfo(ctx, types.TickInfo{
		PoolId:         2,
		TickIndex:      1,
		LiquidityGross: math.LegacyOneDec(),
		LiquidityNet:   math.LegacyOneDec(),
		FeeGrowth:      feeGrowth,
	})

	tickInfo, err := k.GetTickInfo(ctx, 1, 1)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(1))
	require.Equal(t, tickInfo.FeeGrowth.String(), "1.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross.String(), "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "1.000000000000000000")

	tickInfo, err = k.GetTickInfo(ctx, 2, 1)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(2))
	require.Equal(t, tickInfo.TickIndex, int64(1))
	require.Equal(t, tickInfo.FeeGrowth.String(), "1.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross.String(), "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "1.000000000000000000")

	tickInfo, err = k.GetTickInfo(ctx, 1, 2)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(2))
	require.Equal(t, tickInfo.FeeGrowth.String(), "1.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross.String(), "1.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "1.000000000000000000")

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
	k, bk, srv, ctx := setupMsgServer(t)

	// When pool does not exist
	_, err := k.UpsertTick(ctx, 1, 0, math.LegacyNewDec(10), true)
	require.Error(t, err)

	// When pool exist
	wctx := sdk.UnwrapSDKContext(ctx)

	bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	sender := sdk.AccAddress("sender")
	_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
		Authority:  sender.String(),
		DenomBase:  "base",
		DenomQuote: "quote",
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "0.5",
	})
	require.NoError(t, err)

	tickEmpty, err := k.UpsertTick(ctx, 1, 0, math.LegacyNewDec(10), true)
	require.NoError(t, err)
	require.False(t, tickEmpty)

	// Check state change in tickInfo
	tickInfo, err := k.GetTickInfo(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross.String(), "10.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "-10.000000000000000000")

	// Tick's available
	tickEmpty, err = k.UpsertTick(ctx, 1, 0, math.LegacyNewDec(10), false)
	require.NoError(t, err)
	require.False(t, tickEmpty)

	// Check state change in tickInfo
	tickInfo, err = k.GetTickInfo(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross.String(), "20.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "0.000000000000000000")

	// Negative deltaLiquidity
	tickEmpty, err = k.UpsertTick(ctx, 1, 0, math.LegacyNewDec(-20), false)
	require.NoError(t, err)
	require.False(t, tickEmpty)

	// Check state change in tickInfo
	tickInfo, err = k.GetTickInfo(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross.String(), "0.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "-20.000000000000000000")
}

func TestNewTickInfo(t *testing.T) {

	k, bk, srv, ctx := setupMsgServer(t)

	// When pool does not exist
	_, err := k.NewTickInfo(ctx, 1, 0)
	require.Error(t, err)

	// When empty pool exist
	wctx := sdk.UnwrapSDKContext(ctx)

	bk.EXPECT().IsSendEnabledCoins(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	bk.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	sender := sdk.AccAddress("sender")
	_, err = srv.CreatePool(wctx, &types.MsgCreatePool{
		Authority:  sender.String(),
		DenomBase:  "base",
		DenomQuote: "quote",
		FeeRate:    "0.01",
		PriceRatio: "1.0001",
		BaseOffset: "0.5",
	})
	require.NoError(t, err)

	tickInfo, err := k.NewTickInfo(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "")
	require.Equal(t, tickInfo.LiquidityGross.String(), "0.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "0.000000000000000000")

	// When pool accumulator has positive accumulation value
	accumulator, err := k.GetFeeAccumulator(ctx, 1)
	require.NoError(t, err)
	accumulator.AccumValue = sdk.DecCoins{sdk.NewInt64DecCoin("denom", 100)}
	err = k.SetAccumulator(ctx, accumulator)
	require.NoError(t, err)

	tickInfo, err = k.NewTickInfo(ctx, 1, 0)
	require.NoError(t, err)
	require.Equal(t, tickInfo.PoolId, uint64(1))
	require.Equal(t, tickInfo.TickIndex, int64(0))
	require.Equal(t, tickInfo.FeeGrowth.String(), "100.000000000000000000denom")
	require.Equal(t, tickInfo.LiquidityGross.String(), "0.000000000000000000")
	require.Equal(t, tickInfo.LiquidityNet.String(), "0.000000000000000000")
}

// TODO: add test for crossTick
// TODO: add test for ParseTickFromBz
