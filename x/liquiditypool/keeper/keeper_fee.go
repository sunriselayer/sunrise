package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

var emptyCoins = sdk.DecCoins(nil)

func (k Keeper) createFeeAccumulator(ctx context.Context, poolId uint64) error {
	err := k.InitAccumulator(ctx, types.KeyFeePoolAccumulator(poolId))
	if err != nil {
		return errorsmod.Wrapf(err, "failed to create fee accumulator: %d", poolId)
	}
	return nil
}

func (k Keeper) GetFeeAccumulator(ctx context.Context, poolId uint64) (types.AccumulatorObject, error) {
	return k.GetAccumulator(ctx, types.KeyFeePoolAccumulator(poolId))
}

func (k Keeper) SetAccumulatorPositionFeeAccumulator(ctx sdk.Context, poolId uint64, lowerTick, upperTick int64, positionId uint64, liquidityDelta math.LegacyDec) error {
	feeAccumulator, err := k.GetFeeAccumulator(ctx, poolId)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to get fee accumulator: %d", poolId)
	}

	positionKey := types.KeyFeePositionAccumulator(positionId)

	hasPosition := k.HasPosition(ctx, feeAccumulator.Name, positionKey)

	feeGrowthOutside, err := k.getFeeGrowthOutside(ctx, poolId, lowerTick, upperTick)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to get fee growth outside: %d", poolId)
	}

	feeGrowthInside, _ := feeAccumulator.AccumValue.SafeSub(feeGrowthOutside)

	if !hasPosition {
		if !liquidityDelta.IsPositive() {
			return types.ErrNonPositiveLiquidity
		}

		if err := k.NewPositionIntervalAccumulation(ctx, feeAccumulator.Name, positionKey, liquidityDelta, feeGrowthInside); err != nil {
			return errorsmod.Wrapf(err, "failed to new position interval accumulation: %d", poolId)
		}
	} else {
		err = k.updatePositionToInitValuePlusGrowthOutside(ctx, feeAccumulator.Name, positionKey, feeGrowthOutside)
		if err != nil {
			return errorsmod.Wrapf(err, "failed to update position to init value plus growth outside: %d", poolId)
		}

		err = k.UpdatePositionIntervalAccumulation(ctx, feeAccumulator.Name, positionKey, liquidityDelta, feeGrowthInside)
		if err != nil {
			return errorsmod.Wrapf(err, "failed to update position interval accumulation: %d", poolId)
		}
	}

	return nil
}

func (k Keeper) getFeeGrowthOutside(ctx sdk.Context, poolId uint64, lowerTick, upperTick int64) (sdk.DecCoins, error) {
	pool, found, err := k.GetPool(ctx, poolId)
	if err != nil {
		return sdk.DecCoins{}, errorsmod.Wrapf(err, "failed to get pool: %d", poolId)
	}
	if !found {
		return sdk.DecCoins{}, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", poolId)
	}
	currentTick := pool.GetCurrentTick()

	lowerTickInfo, err := k.GetTickInfo(ctx, poolId, lowerTick)
	if err != nil {
		return sdk.DecCoins{}, errorsmod.Wrapf(err, "failed to get lower tick info: %d", lowerTick)
	}
	upperTickInfo, err := k.GetTickInfo(ctx, poolId, upperTick)
	if err != nil {
		return sdk.DecCoins{}, errorsmod.Wrapf(err, "failed to get upper tick info: %d", upperTick)
	}

	poolFeeAccumulator, err := k.GetFeeAccumulator(ctx, poolId)
	if err != nil {
		return sdk.DecCoins{}, errorsmod.Wrapf(err, "failed to get fee accumulator: %d", poolId)
	}
	poolFeeGrowth := poolFeeAccumulator.AccumValue

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
		return feeAccumulator.AccumValue, nil
	}

	return emptyCoins, nil
}

func (k Keeper) collectFees(ctx sdk.Context, sender sdk.AccAddress, positionId uint64) (sdk.Coins, error) {
	position, found, err := k.GetPosition(ctx, positionId)
	if err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(err, "failed to get position: %d", positionId)
	}
	if !found {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrPositionNotFound, "position id: %d", positionId)
	}

	if sender.String() != position.Address {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrNotPositionOwner, "sender: %s, position owner: %s", sender.String(), position.Address)
	}

	feesClaimed, err := k.prepareClaimableFees(ctx, positionId)
	if err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(err, "failed to prepare claimable fees: %d", positionId)
	}

	if feesClaimed.IsZero() {
		return sdk.Coins{}, nil
	}

	pool, found, err := k.GetPool(ctx, position.PoolId)
	if err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(err, "failed to get pool: %d", position.PoolId)
	}
	if !found {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", position.PoolId)
	}
	if err := k.bankKeeper.SendCoins(ctx, pool.GetFeesAddress(), sender, feesClaimed); err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(err, "failed to send coins: %d", position.PoolId)
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventClaimRewards{
		Address:    sender.String(),
		PositionId: positionId,
		Rewards:    feesClaimed,
	}); err != nil {
		return sdk.Coins{}, err
	}

	return feesClaimed, nil
}

func (k Keeper) GetClaimableFees(ctx sdk.Context, positionId uint64) (sdk.Coins, error) {
	cacheCtx, _ := ctx.CacheContext()
	return k.prepareClaimableFees(cacheCtx, positionId)
}

func (k Keeper) prepareClaimableFees(ctx sdk.Context, positionId uint64) (sdk.Coins, error) {
	position, found, err := k.GetPosition(ctx, positionId)
	if err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(err, "failed to get position: %d", positionId)
	}
	if !found {
		return sdk.Coins{}, errorsmod.Wrapf(types.ErrPositionNotFound, "position id: %d", positionId)
	}

	feeAccumulator, err := k.GetFeeAccumulator(ctx, position.PoolId)
	if err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(err, "failed to get fee accumulator: %d", position.PoolId)
	}

	positionKey := types.KeyFeePositionAccumulator(positionId)

	hasPosition := k.HasPosition(ctx, feeAccumulator.Name, positionKey)
	if !hasPosition {
		return nil, types.ErrFeePositionNotFound
	}

	feeGrowthOutside, err := k.getFeeGrowthOutside(ctx, position.PoolId, position.LowerTick, position.UpperTick)
	if err != nil {
		return nil, err
	}

	feesClaimed, forfeitedDust, err := k.updateAccumAndClaimRewards(ctx, feeAccumulator, positionKey, feeGrowthOutside)
	if err != nil {
		return nil, err
	}

	if !forfeitedDust.IsZero() {
		feeAccumulator, err := k.GetFeeAccumulator(ctx, position.PoolId)
		if err != nil {
			return nil, err
		}

		totalSharesRemaining, err := math.LegacyNewDecFromStr(feeAccumulator.TotalShares)
		if err != nil {
			return nil, err
		}
		if !totalSharesRemaining.IsZero() {
			forfeitedDustPerShareScaled := forfeitedDust.QuoDecTruncate(totalSharesRemaining)
			err = k.AddToAccumulator(ctx, feeAccumulator, forfeitedDustPerShareScaled)
			if err != nil {
				return nil, err
			}
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

func (k Keeper) updatePositionToInitValuePlusGrowthOutside(ctx context.Context, accumName, positionKey string, growthOutside sdk.DecCoins) error {
	position, err := k.GetAccumulatorPosition(ctx, accumName, positionKey)
	if err != nil {
		return err
	}

	intervalAccumulationOutside := position.AccumValuePerShare.Add(growthOutside...)
	err = k.SetPositionIntervalAccumulation(ctx, accumName, positionKey, intervalAccumulationOutside)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) updateAccumAndClaimRewards(ctx context.Context, accumulator types.AccumulatorObject, positionKey string, growthOutside sdk.DecCoins) (sdk.Coins, sdk.DecCoins, error) {
	err := k.updatePositionToInitValuePlusGrowthOutside(ctx, accumulator.Name, positionKey, growthOutside)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, err
	}

	incentivesClaimedCurrAccum, dust, err := k.ClaimRewards(ctx, accumulator.Name, positionKey)
	if err != nil {
		return sdk.Coins{}, sdk.DecCoins{}, err
	}

	hasPosition := k.HasPosition(ctx, accumulator.Name, positionKey)
	if hasPosition {
		currentGrowthInsideForPosition, _ := accumulator.AccumValue.SafeSub(growthOutside)
		err := k.SetPositionIntervalAccumulation(ctx, accumulator.Name, positionKey, currentGrowthInsideForPosition)
		if err != nil {
			return sdk.Coins{}, sdk.DecCoins{}, err
		}
	}

	return incentivesClaimedCurrAccum, dust, nil
}
