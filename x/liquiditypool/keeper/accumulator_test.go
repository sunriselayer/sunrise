package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/sunriselayer/sunrise/testutil/keeper"
)

// TODO: add test for NewPosition
// TODO: add test for NewPositionIntervalAccumulation
// TODO: add test for AddToPosition
// TODO: add test for AddToPositionIntervalAccumulation
// TODO: add test for RemoveFromPosition
// TODO: add test for RemoveFromPositionIntervalAccumulation
// TODO: add test for UpdatePositionIntervalAccumulation
// TODO: add test for SetPositionIntervalAccumulation
// TODO: add test for DeletePosition
// TODO: add test for deletePosition
// TODO: add test for GetAccumulatorPositionSize
// TODO: add test for ClaimRewards
// TODO: add test for AddToUnclaimedRewards
// TODO: add test for GetTotalRewards

func TestAccumulatorStore(t *testing.T) {
	k, _, ctx := keepertest.LiquiditypoolKeeper(t)

	err := k.InitAccumulator(ctx, "accumulator1")
	require.NoError(t, err)
	err = k.InitAccumulator(ctx, "accumulator2")
	require.NoError(t, err)

	accumulator, err := k.GetAccumulator(ctx, "accumulator1")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator1")
	require.Equal(t, accumulator.AccumValue.String(), "")
	require.Equal(t, accumulator.TotalShares.String(), "0.000000000000000000")

	accumulator, err = k.GetAccumulator(ctx, "accumulator2")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator2")
	require.Equal(t, accumulator.AccumValue.String(), "")
	require.Equal(t, accumulator.TotalShares.String(), "0.000000000000000000")

	k.AddToAccumulator(ctx, accumulator, sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(100))))
	accumulator, err = k.GetAccumulator(ctx, "accumulator2")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator2")
	require.Equal(t, accumulator.AccumValue.String(), "100.000000000000000000denom")
	require.Equal(t, accumulator.TotalShares.String(), "0.000000000000000000")

	accumulator.AccumValue = sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	err = k.SetAccumulator(ctx, accumulator)
	require.NoError(t, err)
	accumulator, err = k.GetAccumulator(ctx, "accumulator2")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator2")
	require.Equal(t, accumulator.AccumValue.String(), "1.000000000000000000denom")

	_, err = k.GetAccumulator(ctx, "accumulator3")
	require.Error(t, err)

	accumulators := k.GetAllAccumulators(ctx)
	require.Len(t, accumulators, 2)
}

func TestAccumulatorPositionStore(t *testing.T) {
	k, _, ctx := keepertest.LiquiditypoolKeeper(t)

	// Get not available position
	_, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.Error(t, err)
	hasPosition := k.HasPosition(ctx, "accumulator", "index")
	require.False(t, hasPosition)

	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	unclaimedRewardsTotal := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(2)))
	k.SetAccumulatorPosition(ctx, "accumulator", accmulatorValuePerShare, "index", math.LegacyOneDec(), unclaimedRewardsTotal)
	position, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares.String(), "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "2.000000000000000000denom")

	k.SetAccumulatorPosition(ctx, "accumulator", accmulatorValuePerShare, "index2", math.LegacyOneDec(), unclaimedRewardsTotal)
	position, err = k.GetAccumulatorPosition(ctx, "accumulator", "index2")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index2")
	require.Equal(t, position.NumShares.String(), "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "2.000000000000000000denom")

	k.SetAccumulatorPosition(ctx, "accumulator2", accmulatorValuePerShare, "index", math.LegacyOneDec(), unclaimedRewardsTotal)
	position, err = k.GetAccumulatorPosition(ctx, "accumulator2", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator2")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares.String(), "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "2.000000000000000000denom")

	hasPosition = k.HasPosition(ctx, "accumulator2", "index")
	require.True(t, hasPosition)

	positions := k.GetAllAccumulatorPositions(ctx)
	require.Len(t, positions, 3)
}
