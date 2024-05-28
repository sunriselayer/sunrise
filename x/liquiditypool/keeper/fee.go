package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

var emptyCoins = sdk.DecCoins(nil)

func (k Keeper) createFeeAccumulator(ctx context.Context, poolId uint64) error {
	err := MakeAccumulator(k.storeService.OpenKVStore(ctx), types.KeyFeePoolAccumulator(poolId))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetFeeAccumulator(ctx context.Context, poolId uint64) (*AccumulatorObject, error) {
	return GetAccumulator(k.storeService.OpenKVStore(ctx), types.KeyFeePoolAccumulator(poolId))
}

func (k Keeper) initOrUpdatePositionFeeAccumulator(ctx sdk.Context, poolId uint64, lowerTick, upperTick int64, positionId uint64, liquidityDelta math.LegacyDec) error {
	feeAccumulator, err := k.GetFeeAccumulator(ctx, poolId)
	if err != nil {
		return err
	}

	positionKey := types.KeyFeePositionAccumulator(positionId)

	hasPosition := feeAccumulator.HasPosition(positionKey)

	feeGrowthOutside, err := k.getFeeGrowthOutside(ctx, poolId, lowerTick, upperTick)
	if err != nil {
		return err
	}

	feeGrowthInside, _ := feeAccumulator.GetValue().SafeSub(feeGrowthOutside)

	if !hasPosition {
		if !liquidityDelta.IsPositive() {
			return types.ErrNonPositiveLiquidity
		}

		if err := feeAccumulator.NewPositionIntervalAccumulation(positionKey, liquidityDelta, feeGrowthInside); err != nil {
			return err
		}
	} else {
		err = updatePositionToInitValuePlusGrowthOutside(feeAccumulator, positionKey, feeGrowthOutside)
		if err != nil {
			return err
		}

		err = feeAccumulator.UpdatePositionIntervalAccumulation(positionKey, liquidityDelta, feeGrowthInside)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) getFeeGrowthOutside(ctx sdk.Context, poolId uint64, lowerTick, upperTick int64) (sdk.DecCoins, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.DecCoins{}, types.ErrPoolNotFound
	}
	currentTick := pool.GetCurrentTick()

	lowerTickInfo, err := k.GetTickInfo(ctx, poolId, lowerTick)
	if err != nil {
		return sdk.DecCoins{}, err
	}
	upperTickInfo, err := k.GetTickInfo(ctx, poolId, upperTick)
	if err != nil {
		return sdk.DecCoins{}, err
	}

	poolFeeAccumulator, err := k.GetFeeAccumulator(ctx, poolId)
	if err != nil {
		return sdk.DecCoins{}, err
	}
	poolFeeGrowth := poolFeeAccumulator.GetValue()

	feeGrowthAboveUpperTick := calculateFeeGrowth(upperTick, upperTickInfo.FeeGrowth, currentTick, poolFeeGrowth, true)
	feeGrowthBelowLowerTick := calculateFeeGrowth(lowerTick, lowerTickInfo.FeeGrowth, currentTick, poolFeeGrowth, false)

	return feeGrowthAboveUpperTick.Add(feeGrowthBelowLowerTick...), nil
}

func (k Keeper) getInitialFeeGrowth(ctx context.Context, pool types.Pool, tick int64) (sdk.DecCoins, error) {
	currentTick := pool.GetCurrentTick()
	if currentTick >= tick {
		feeAccumulator, err := k.GetFeeAccumulator(ctx, pool.GetId())
		if err != nil {
			return sdk.DecCoins{}, err
		}
		return feeAccumulator.GetValue(), nil
	}

	return emptyCoins, nil
}

func (k Keeper) collectFees(ctx sdk.Context, sender sdk.AccAddress, positionId uint64) (sdk.Coins, error) {
	position, found := k.GetPosition(ctx, positionId)
	if !found {
		return sdk.Coins{}, types.ErrPositionNotFound
	}

	if sender.String() != position.Address {
		return sdk.Coins{}, types.ErrNotPositionOwner
	}

	feesClaimed, err := k.prepareClaimableFees(ctx, positionId)
	if err != nil {
		return sdk.Coins{}, err
	}

	if feesClaimed.IsZero() {
		return sdk.Coins{}, nil
	}

	pool, found := k.GetPool(ctx, position.PoolId)
	if !found {
		return sdk.Coins{}, types.ErrPoolNotFound
	}
	if err := k.bankKeeper.SendCoins(ctx, pool.GetFeesAddress(), sender, feesClaimed); err != nil {
		return sdk.Coins{}, err
	}

	return feesClaimed, nil
}

func (k Keeper) GetClaimableFees(ctx sdk.Context, positionId uint64) (sdk.Coins, error) {
	cacheCtx, _ := ctx.CacheContext()
	return k.prepareClaimableFees(cacheCtx, positionId)
}

func (k Keeper) prepareClaimableFees(ctx sdk.Context, positionId uint64) (sdk.Coins, error) {
	position, found := k.GetPosition(ctx, positionId)
	if !found {
		return nil, types.ErrPositionNotFound
	}

	feeAccumulator, err := k.GetFeeAccumulator(ctx, position.PoolId)
	if err != nil {
		return nil, err
	}

	positionKey := types.KeyFeePositionAccumulator(positionId)

	hasPosition := feeAccumulator.HasPosition(positionKey)
	if !hasPosition {
		return nil, types.ErrFeePositionNotFound
	}

	feeGrowthOutside, err := k.getFeeGrowthOutside(ctx, position.PoolId, position.LowerTick, position.UpperTick)
	if err != nil {
		return nil, err
	}

	feesClaimedScaled, forfeitedDustScaled, err := updateAccumAndClaimRewards(feeAccumulator, positionKey, feeGrowthOutside)
	if err != nil {
		return nil, err
	}

	spreadFactorScalingFactor, err := k.getSpreadFactorScalingFactorForPool(ctx, position.PoolId)
	if err != nil {
		return nil, err
	}

	feesClaimed := sdk.NewCoins()
	forfeitedDust := sdk.DecCoins{}
	if spreadFactorScalingFactor.Equal(math.LegacyOneDec()) {
		feesClaimed = feesClaimedScaled
		forfeitedDust = forfeitedDustScaled
	} else {
		for _, coin := range feesClaimedScaled {
			scaledCoinAmt := scaleDownFeeAmount(coin.Amount, spreadFactorScalingFactor)
			if !scaledCoinAmt.IsZero() {
				feesClaimed = append(feesClaimed, sdk.NewCoin(coin.Denom, scaledCoinAmt))
			}
		}
	}

	if !forfeitedDust.IsZero() {
		feeAccumulator, err := k.GetFeeAccumulator(ctx, position.PoolId)
		if err != nil {
			return nil, err
		}

		totalSharesRemaining := feeAccumulator.GetTotalShares()

		if !totalSharesRemaining.IsZero() {
			forfeitedDustPerShareScaled := forfeitedDust.QuoDecTruncate(totalSharesRemaining)
			feeAccumulator.AddToAccumulator(forfeitedDustPerShareScaled)
		}
	}

	return feesClaimed, nil
}

func calculateFeeGrowth(targetTick int64, ticksFeeGrowthOppositeDirectionOfLastTraversal sdk.DecCoins, currentTick int64, feesGrowthGlobal sdk.DecCoins, isUpperTick bool) sdk.DecCoins {
	if (isUpperTick && currentTick >= targetTick) || (!isUpperTick && currentTick < targetTick) {
		return feesGrowthGlobal.Sub(ticksFeeGrowthOppositeDirectionOfLastTraversal)
	}
	return ticksFeeGrowthOppositeDirectionOfLastTraversal
}

func updatePositionToInitValuePlusGrowthOutside(accumulator *AccumulatorObject, positionKey string, growthOutside sdk.DecCoins) error {
	position, err := GetPosition(accumulator, positionKey)
	if err != nil {
		return err
	}

	intervalAccumulationOutside := position.AccumValuePerShare.Add(growthOutside...)
	err = accumulator.SetPositionIntervalAccumulation(positionKey, intervalAccumulationOutside)
	if err != nil {
		return err
	}
	return nil
}

func scaleDownFeeAmount(incentiveAmount math.Int, scalingFactor math.LegacyDec) (scaledTotalEmittedAmount math.Int) {
	return incentiveAmount.ToLegacyDec().QuoTruncateMut(scalingFactor).TruncateInt()
}

func (k Keeper) getSpreadFactorScalingFactorForPool(ctx sdk.Context, poolID uint64) (math.LegacyDec, error) {
	return math.LegacyOneDec(), nil
}

func updateAccumAndClaimRewards(accum *AccumulatorObject, positionKey string, growthOutside sdk.DecCoins) (sdk.Coins, sdk.DecCoins, error) {
	err := updatePositionToInitValuePlusGrowthOutside(accum, positionKey, growthOutside)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, err
	}

	incentivesClaimedCurrAccum, dust, err := accum.ClaimRewards(positionKey)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, err
	}

	hasPosition := accum.HasPosition(positionKey)

	if hasPosition {
		currentGrowthInsideForPosition, _ := accum.GetValue().SafeSub(growthOutside)
		err := accum.SetPositionIntervalAccumulation(positionKey, currentGrowthInsideForPosition)
		if err != nil {
			return sdk.Coins{}, sdk.DecCoins{}, err
		}
	}

	return incentivesClaimedCurrAccum, dust, nil
}
