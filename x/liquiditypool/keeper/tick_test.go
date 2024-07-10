package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// TODO: add test for initOrUpdateTick
// TODO: add test for crossTick
// TODO: add test for newTickInfo
// TODO: add test for ParseTickFromBz

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
