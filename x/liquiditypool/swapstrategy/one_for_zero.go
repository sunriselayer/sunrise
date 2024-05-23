package swapstrategy

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	"cosmossdk.io/core/store"
)

type oneForZeroStrategy struct {
	sqrtPriceLimit math.LegacyDec
	storeService   store.KVStoreService
	spreadFactor   math.LegacyDec

	oneMinusSpreadFactor math.LegacyDec
	spfOverOneMinusSpf   math.LegacyDec
}

var _ SwapStrategy = (*oneForZeroStrategy)(nil)

func (s oneForZeroStrategy) ZeroForOne() bool { return false }

func (s oneForZeroStrategy) GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec {
	if nextTickSqrtPrice.GT(s.sqrtPriceLimit) {
		return s.sqrtPriceLimit
	}
	return nextTickSqrtPrice
}

func (s oneForZeroStrategy) ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountOneInRemaining math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountOneIn := types.CalcAmountQuoteDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, true)

	oneMinusTakerFee := s.getOneMinusSpreadFactor()
	amountOneInRemainingLessSpreadReward := amountOneInRemaining.Mul(oneMinusTakerFee)

	var sqrtPriceNext math.LegacyDec
	if amountOneInRemainingLessSpreadReward.GTE(amountOneIn) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountQuoteInRoundingDown(sqrtPriceCurrent, liquidity, amountOneInRemainingLessSpreadReward)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountOneIn = types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true) // N.B.: if this is false, causes infinite loop
	}

	amountZeroOut := types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)

	spreadRewardChargeTotal := computeSpreadRewardChargePerSwapStepOutGivenIn(hasReachedTarget, amountOneIn, amountOneInRemaining, s.spreadFactor, s.getSpfOverOneMinusSpf)
	return sqrtPriceNext, amountOneIn, amountZeroOut, spreadRewardChargeTotal
}

func (s oneForZeroStrategy) ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountZeroRemainingOut math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountZeroOut := types.CalcAmountBaseDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, false)

	var sqrtPriceNext math.LegacyDec
	if amountZeroRemainingOut.GTE(amountZeroOut) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountBaseOutRoundingUp(sqrtPriceCurrent, liquidity, amountZeroRemainingOut)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountZeroOut = types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)
	}

	amountOneIn := types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)

	spreadRewardChargeTotal := computeSpreadRewardChargeFromAmountIn(amountOneIn, s.getSpfOverOneMinusSpf())

	if amountZeroOut.GT(amountZeroRemainingOut) {
		amountZeroOut = amountZeroRemainingOut
	}

	return sqrtPriceNext, amountZeroOut, amountOneIn, spreadRewardChargeTotal
}

func (s oneForZeroStrategy) getOneMinusSpreadFactor() math.LegacyDec {
	if s.oneMinusSpreadFactor.IsNil() {
		s.oneMinusSpreadFactor = oneDec.Sub(s.spreadFactor)
	}
	return s.oneMinusSpreadFactor
}

func (s oneForZeroStrategy) getSpfOverOneMinusSpf() math.LegacyDec {
	if s.spfOverOneMinusSpf.IsNil() {
		s.spfOverOneMinusSpf = s.spreadFactor.QuoRoundUp(s.getOneMinusSpreadFactor())
	}
	return s.spfOverOneMinusSpf
}

func (s oneForZeroStrategy) InitializeNextTickIterator(ctx sdk.Context, poolId uint64, currentTickIndex int64) dbm.Iterator {
	storeAdapter := runtime.KVStoreAdapter(s.storeService.OpenKVStore(ctx))
	prefixBz := types.KeyTickPrefixByPoolId(poolId)
	prefixStore := prefix.NewStore(storeAdapter, prefixBz)
	startKey := types.TickIndexToBytes(currentTickIndex)
	iter := prefixStore.Iterator(startKey, nil)

	for ; iter.Valid(); iter.Next() {
		tick, err := types.TickIndexFromBytes(iter.Key())
		if err != nil {
			iter.Close()
			panic(fmt.Errorf("invalid tick index (%s): %v", string(iter.Key()), err))
		}

		if tick > currentTickIndex {
			break
		}
	}
	return iter
}

func (s oneForZeroStrategy) SetLiquidityDeltaSign(deltaLiquidity math.LegacyDec) math.LegacyDec {
	return deltaLiquidity
}

func (s oneForZeroStrategy) UpdateTickAfterCrossing(nextTick int64) int64 {
	return nextTick
}

func (s oneForZeroStrategy) ValidateSqrtPrice(sqrtPrice math.LegacyDec, currentSqrtPrice math.LegacyDec) error {
	if sqrtPrice.LT(currentSqrtPrice) || sqrtPrice.GT(types.MaxSqrtPrice) {
		return types.ErrInvalidSqrtPrice
	}
	return nil
}
