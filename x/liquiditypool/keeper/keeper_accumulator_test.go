package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func TestAccumulatorStore(t *testing.T) {
	f := initFixture(t)
	ctx := f.ctx
	k := f.keeper

	err := k.InitAccumulator(ctx, "accumulator1")
	require.NoError(t, err)
	err = k.InitAccumulator(ctx, "accumulator2")
	require.NoError(t, err)

	accumulator, err := k.GetAccumulator(ctx, "accumulator1")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator1")
	require.Equal(t, accumulator.AccumValue.String(), "")
	require.Equal(t, accumulator.TotalShares, "0.000000000000000000")

	accumulator, err = k.GetAccumulator(ctx, "accumulator2")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator2")
	require.Equal(t, accumulator.AccumValue.String(), "")
	require.Equal(t, accumulator.TotalShares, "0.000000000000000000")

	err = k.AddToAccumulator(ctx, accumulator, sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(100))))
	require.NoError(t, err)
	accumulator, err = k.GetAccumulator(ctx, "accumulator2")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator2")
	require.Equal(t, accumulator.AccumValue.String(), "100.000000000000000000denom")
	require.Equal(t, accumulator.TotalShares, "0.000000000000000000")

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
	f := initFixture(t)
	ctx := f.ctx
	k := f.keeper

	// Get not available position
	_, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.Error(t, err)
	hasPosition := k.HasPosition(ctx, "accumulator", "index")
	require.False(t, hasPosition)
	_, err = k.GetAccumulatorPositionSize(ctx, "accumulator2", "index")
	require.Error(t, err)

	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	unclaimedRewardsTotal := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(2)))
	err = k.SetAccumulatorPosition(ctx, "accumulator", accmulatorValuePerShare, "index", math.LegacyOneDec(), unclaimedRewardsTotal)
	require.NoError(t, err)
	position, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares, "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "2.000000000000000000denom")

	err = k.SetAccumulatorPosition(ctx, "accumulator", accmulatorValuePerShare, "index2", math.LegacyOneDec(), unclaimedRewardsTotal)
	require.NoError(t, err)
	position, err = k.GetAccumulatorPosition(ctx, "accumulator", "index2")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index2")
	require.Equal(t, position.NumShares, "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "2.000000000000000000denom")

	err = k.SetAccumulatorPosition(ctx, "accumulator2", accmulatorValuePerShare, "index", math.LegacyOneDec(), unclaimedRewardsTotal)
	require.NoError(t, err)
	position, err = k.GetAccumulatorPosition(ctx, "accumulator2", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator2")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares, "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "2.000000000000000000denom")

	positionSize, err := k.GetAccumulatorPositionSize(ctx, "accumulator2", "index")
	require.NoError(t, err)
	require.Equal(t, positionSize.String(), "1.000000000000000000")

	hasPosition = k.HasPosition(ctx, "accumulator2", "index")
	require.True(t, hasPosition)

	positions := k.GetAllAccumulatorPositions(ctx)
	require.Len(t, positions, 3)
}

func TestNewPositionIntervalAccumulation(t *testing.T) {
	f := initFixture(t)
	ctx := f.ctx
	k := f.keeper
	// when accumulator does not exist
	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	err := k.NewPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.Error(t, err)

	// when accumulator exists
	err = k.InitAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	err = k.NewPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.NoError(t, err)

	// check accumulator change
	accumulator, err := k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator")
	require.Equal(t, accumulator.AccumValue.String(), "")
	require.Equal(t, accumulator.TotalShares, "1.000000000000000000")

	// check accumulator position change
	position, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares, "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "")
}

func TestAddToPositionIntervalAccumulation(t *testing.T) {
	f := initFixture(t)
	ctx := f.ctx
	k := f.keeper
	// when new shares is negative
	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	err := k.AddToPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec().Neg(), accmulatorValuePerShare)
	require.Error(t, err)

	// when position does not exist
	err = k.AddToPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.Error(t, err)

	// when accumulator and position exists
	err = k.InitAccumulator(ctx, "accumulator")
	require.NoError(t, err)

	accumulator, err := k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	accumulator.AccumValue = accumulator.AccumValue.Add(accmulatorValuePerShare...).Add(accmulatorValuePerShare...)
	err = k.SetAccumulator(ctx, accumulator)
	require.NoError(t, err)

	err = k.NewPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.NoError(t, err)
	err = k.AddToPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.NoError(t, err)

	// check accumulator change
	accumulator, err = k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator")
	require.Equal(t, accumulator.AccumValue.String(), "2.000000000000000000denom")
	require.Equal(t, accumulator.TotalShares, "2.000000000000000000")

	// check accumulator position change
	position, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares, "2.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "1.000000000000000000denom")
}

func TestRemoveFromPositionIntervalAccumulation(t *testing.T) {
	f := initFixture(t)
	ctx := f.ctx
	k := f.keeper
	// when new shares is negative
	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	err := k.RemoveFromPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec().Neg(), accmulatorValuePerShare)
	require.Error(t, err)

	// when position does not exist
	err = k.RemoveFromPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.Error(t, err)

	// when accumulator and position exists
	err = k.InitAccumulator(ctx, "accumulator")
	require.NoError(t, err)

	accumulator, err := k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	accumulator.AccumValue = accumulator.AccumValue.Add(accmulatorValuePerShare...).Add(accmulatorValuePerShare...)
	err = k.SetAccumulator(ctx, accumulator)
	require.NoError(t, err)

	err = k.NewPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.NoError(t, err)
	err = k.RemoveFromPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.NoError(t, err)

	// check accumulator change
	accumulator, err = k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator")
	require.Equal(t, accumulator.AccumValue.String(), "2.000000000000000000denom")
	require.Equal(t, accumulator.TotalShares, "0.000000000000000000")

	// check accumulator position change
	position, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares, "0.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "1.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "1.000000000000000000denom")
}

func TestGetTotalRewards(t *testing.T) {
	// When accumulator value is lower than position value
	oneDecCoins := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	twoDecCoins := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(2)))
	emptyDecCoins := sdk.NewDecCoins()
	rewards := keeper.GetTotalRewards(types.AccumulatorObject{
		AccumValue: oneDecCoins,
	}, types.AccumulatorPosition{
		AccumValuePerShare:    twoDecCoins,
		NumShares:             "1.000000000000000000",
		UnclaimedRewardsTotal: emptyDecCoins,
	})
	require.Equal(t, rewards.String(), "")

	// When accumulator value is equal to position value
	rewards = keeper.GetTotalRewards(types.AccumulatorObject{
		AccumValue: oneDecCoins,
	}, types.AccumulatorPosition{
		AccumValuePerShare:    oneDecCoins,
		NumShares:             "1.000000000000000000",
		UnclaimedRewardsTotal: emptyDecCoins,
	})
	require.Equal(t, rewards.String(), "")

	// When accumulator value is greater than position value
	rewards = keeper.GetTotalRewards(types.AccumulatorObject{
		AccumValue: twoDecCoins,
	}, types.AccumulatorPosition{
		AccumValuePerShare:    oneDecCoins,
		NumShares:             "1.000000000000000000",
		UnclaimedRewardsTotal: emptyDecCoins,
	})
	require.Equal(t, rewards.String(), "1.000000000000000000denom")

	// When position numShares is zero
	rewards = keeper.GetTotalRewards(types.AccumulatorObject{
		AccumValue: twoDecCoins,
	}, types.AccumulatorPosition{
		AccumValuePerShare:    oneDecCoins,
		NumShares:             "0.000000000000000000",
		UnclaimedRewardsTotal: emptyDecCoins,
	})
	require.Equal(t, rewards.String(), "")

	// When position numShares is negative
	rewards = keeper.GetTotalRewards(types.AccumulatorObject{
		AccumValue: twoDecCoins,
	}, types.AccumulatorPosition{
		AccumValuePerShare:    oneDecCoins,
		NumShares:             "-1.000000000000000000",
		UnclaimedRewardsTotal: emptyDecCoins,
	})
	require.Equal(t, rewards.String(), "")
}

// func TestLiquidateAndDeletePosition(t *testing.T) {
// 	k, _, ctx := keepertest.LiquiditypoolKeeper(t)
// 	// when position does not exist
// 	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
// 	_, err := k.LiquidateAndDeletePosition(ctx, "accumulator", "index")
// 	require.Error(t, err)

// 	// when accumulator and position exists
// 	err = k.InitAccumulator(ctx, "accumulator")
// 	require.NoError(t, err)

// 	accumulator, err := k.GetAccumulator(ctx, "accumulator")
// 	require.NoError(t, err)
// 	accumulator.AccumValue = accumulator.AccumValue.Add(accmulatorValuePerShare...).Add(accmulatorValuePerShare...)
// 	err = k.SetAccumulator(ctx, accumulator)
// 	require.NoError(t, err)

// 	err = k.NewPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
// 	require.NoError(t, err)
// 	rewards, err := k.LiquidateAndDeletePosition(ctx, "accumulator", "index")
// 	require.NoError(t, err)
// 	require.Equal(t, rewards.String(), "1.000000000000000000denom")

// 	// check accumulator change
// 	accumulator, err = k.GetAccumulator(ctx, "accumulator")
// 	require.NoError(t, err)
// 	require.Equal(t, accumulator.Name, "accumulator")
// 	require.Equal(t, accumulator.AccumValue.String(), "2.000000000000000000denom")
// 	require.Equal(t, accumulator.TotalShares, "0.000000000000000000")

// 	// check accumulator position change
// 	_, err = k.GetAccumulatorPosition(ctx, "accumulator", "index")
// 	require.Error(t, err)
// }

func TestClaimRewards(t *testing.T) {
	f := initFixture(t)
	ctx := f.ctx
	k := f.keeper
	// when new shares is negative
	accmulatorValuePerShare := sdk.NewDecCoins(sdk.NewDecCoin("denom", math.NewInt(1)))
	_, _, err := k.ClaimRewards(ctx, "accumulator", "index")
	require.Error(t, err)

	// when position does not exist
	_, _, err = k.ClaimRewards(ctx, "accumulator", "index")
	require.Error(t, err)

	// when accumulator and position exists
	err = k.InitAccumulator(ctx, "accumulator")
	require.NoError(t, err)

	accumulator, err := k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	accumulator.AccumValue = accumulator.AccumValue.Add(accmulatorValuePerShare...).Add(accmulatorValuePerShare...)
	err = k.SetAccumulator(ctx, accumulator)
	require.NoError(t, err)

	err = k.NewPositionIntervalAccumulation(ctx, "accumulator", "index", math.LegacyOneDec(), accmulatorValuePerShare)
	require.NoError(t, err)
	rewards, dust, err := k.ClaimRewards(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, rewards.String(), "1denom")
	require.Equal(t, dust.String(), "")

	// check accumulator change
	accumulator, err = k.GetAccumulator(ctx, "accumulator")
	require.NoError(t, err)
	require.Equal(t, accumulator.Name, "accumulator")
	require.Equal(t, accumulator.AccumValue.String(), "2.000000000000000000denom")
	require.Equal(t, accumulator.TotalShares, "1.000000000000000000")

	// check accumulator position change
	position, err := k.GetAccumulatorPosition(ctx, "accumulator", "index")
	require.NoError(t, err)
	require.Equal(t, position.Name, "accumulator")
	require.Equal(t, position.Index, "index")
	require.Equal(t, position.NumShares, "1.000000000000000000")
	require.Equal(t, position.AccumValuePerShare.String(), "2.000000000000000000denom")
	require.Equal(t, position.UnclaimedRewardsTotal.String(), "")
}
