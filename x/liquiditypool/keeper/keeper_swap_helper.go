package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

type SwapHelper interface {
	GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec
	ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountRemainingIn math.LegacyDec) (sqrtPriceNext math.LegacyDec, amountInConsumed, amountOutComputed, feeChargeTotal math.LegacyDec)
	ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountRemainingOut math.LegacyDec) (sqrtPriceNext math.LegacyDec, amountOutConsumed, amountInComputed, feeChargeTotal math.LegacyDec)
	NextTickIterator(ctx sdk.Context, storeService store.KVStoreService, poolId uint64, tickIndex int64) dbm.Iterator
	GetLiquidityDeltaSign(liquidityDelta math.LegacyDec) math.LegacyDec
	NextTickAfterCrossing(nextTick int64) (updatedNextTick int64)
	ValidateSqrtPrice(sqrtPriceLimit math.LegacyDec, currentSqrtPrice math.LegacyDec) error
	BaseForQuote() bool
}

func New(baseForQuote bool, sqrtPriceLimit math.LegacyDec, feeRate math.LegacyDec) SwapHelper {
	if baseForQuote {
		return &baseForQuoteHelper{sqrtPriceLimit: sqrtPriceLimit, feeRate: feeRate}
	}
	return &quoteForBaseHelper{sqrtPriceLimit: sqrtPriceLimit, feeRate: feeRate}
}

func GetMultipliedPriceLimit(baseForQuote bool) math.LegacyDec {
	if baseForQuote {
		return types.MinMultipliedSpotPrice
	}
	return types.MaxMultipliedSpotPrice
}

func GetSqrtPriceLimit(multipliedPriceLimit math.LegacyDec, baseForQuote bool) (math.LegacyDec, error) {
	if multipliedPriceLimit.IsZero() {
		if baseForQuote {
			return types.MinSqrtPrice, nil
		}
		return types.MaxSqrtPrice, nil
	}

	if multipliedPriceLimit.LT(types.MinMultipliedSpotPrice) || multipliedPriceLimit.GT(types.MaxMultipliedSpotPrice) {
		return math.LegacyDec{}, types.ErrPriceOutOfBound
	}

	sqrtPriceLimitMultiplied, err := multipliedPriceLimit.ApproxSqrt()
	if err != nil {
		return math.LegacyDec{}, err
	}

	sqrtPriceLimit := sqrtPriceLimitMultiplied.Quo(types.MultiplierSqrt)
	return sqrtPriceLimit, nil
}

func getFeeRateOverOneMinusFeeRate(feeRate math.LegacyDec) math.LegacyDec {
	return feeRate.QuoRoundUp(math.LegacyOneDec().Sub(feeRate))
}

func computeFeeChargeFromInAmount(amountIn math.LegacyDec, feeRateOveroneMinusFeeRate math.LegacyDec) math.LegacyDec {
	return amountIn.MulRoundUp(feeRateOveroneMinusFeeRate)
}

func computeFeeChargePerSwapStepOutGivenIn(hasReachedTarget bool, amountIn, amountSpecifiedRemaining, feeRate math.LegacyDec) math.LegacyDec {
	if feeRate.IsZero() {
		return math.LegacyZeroDec()
	} else if feeRate.IsNegative() {
		panic(fmt.Errorf("fee rate must be non-negative, was (%s)", feeRate))
	}

	var feeChargeTotal math.LegacyDec
	if hasReachedTarget {
		feeChargeTotal = computeFeeChargeFromInAmount(amountIn, getFeeRateOverOneMinusFeeRate(feeRate))
	} else {
		feeChargeTotal = amountSpecifiedRemaining.Sub(amountIn)
	}

	if feeChargeTotal.IsNegative() {
		panic(fmt.Errorf("fee rate charge must be non-negative, was (%s)", feeChargeTotal))
	}

	return feeChargeTotal
}

type baseForQuoteHelper struct {
	sqrtPriceLimit math.LegacyDec
	feeRate        math.LegacyDec
}

var _ SwapHelper = (*baseForQuoteHelper)(nil)

func (s baseForQuoteHelper) BaseForQuote() bool { return true }

func (s baseForQuoteHelper) GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec {
	if nextTickSqrtPrice.LT(s.sqrtPriceLimit) {
		return s.sqrtPriceLimit
	}
	return nextTickSqrtPrice
}

func (s baseForQuoteHelper) ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountBaseInRemaining math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountBaseIn := types.CalcAmountBaseDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, true)
	amountBaseInAfterFee := amountBaseInRemaining.Mul(math.LegacyOneDec().Sub(s.feeRate))

	var sqrtPriceNext math.LegacyDec
	if amountBaseInAfterFee.GTE(amountBaseIn) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountBaseInRoundingUp(sqrtPriceCurrent, liquidity, amountBaseInAfterFee)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)
	if !hasReachedTarget {
		amountBaseIn = types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)
	}

	amountQuoteOut := types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)

	feeChargeTotal := computeFeeChargePerSwapStepOutGivenIn(hasReachedTarget, amountBaseIn, amountBaseInRemaining, s.feeRate)
	return sqrtPriceNext, amountBaseIn, amountQuoteOut, feeChargeTotal
}

func (s baseForQuoteHelper) ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountQuoteRemainingOut math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountQuoteOut := types.CalcAmountQuoteDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, false)

	var sqrtPriceNext math.LegacyDec
	if amountQuoteRemainingOut.GTE(amountQuoteOut) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountQuoteOutRoundingDown(sqrtPriceCurrent, liquidity, amountQuoteRemainingOut)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountQuoteOut = types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)
	}

	amountBaseIn := types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)

	feeChargeTotal := computeFeeChargeFromInAmount(amountBaseIn, getFeeRateOverOneMinusFeeRate(s.feeRate))

	if amountQuoteOut.GT(amountQuoteRemainingOut) {
		amountQuoteOut = amountQuoteRemainingOut
	}

	return sqrtPriceNext, amountQuoteOut, amountBaseIn, feeChargeTotal
}

func (s baseForQuoteHelper) NextTickIterator(ctx sdk.Context, storeService store.KVStoreService, poolId uint64, currentTickIndex int64) dbm.Iterator {
	storeAdapter := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	prefixBz := types.KeyTickPrefixByPoolId(poolId)
	prefixStore := prefix.NewStore(storeAdapter, prefixBz)
	startKey := types.TickIndexToBytes(currentTickIndex + 1)
	iter := prefixStore.ReverseIterator(nil, startKey)

	for ; iter.Valid(); iter.Next() {
		tick, err := types.TickIndexFromBytes(iter.Key())
		if err != nil {
			iter.Close()
			panic(fmt.Errorf("invalid tick index (%s): %v", string(iter.Key()), err))
		}
		if tick <= currentTickIndex {
			break
		}
	}
	return iter
}

func (s baseForQuoteHelper) GetLiquidityDeltaSign(deltaLiquidity math.LegacyDec) math.LegacyDec {
	return deltaLiquidity.Neg()
}

func (s baseForQuoteHelper) NextTickAfterCrossing(nextTick int64) int64 {
	return nextTick - 1
}

func (s baseForQuoteHelper) ValidateSqrtPrice(sqrtPrice math.LegacyDec, currentSqrtPrice math.LegacyDec) error {
	if sqrtPrice.GT(currentSqrtPrice) || sqrtPrice.LT(types.MinSqrtPrice) {
		return types.ErrInvalidSqrtPrice
	}
	return nil
}

type quoteForBaseHelper struct {
	sqrtPriceLimit math.LegacyDec
	feeRate        math.LegacyDec
}

var _ SwapHelper = (*quoteForBaseHelper)(nil)

func (s quoteForBaseHelper) BaseForQuote() bool { return false }

func (s quoteForBaseHelper) GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec {
	if nextTickSqrtPrice.GT(s.sqrtPriceLimit) {
		return s.sqrtPriceLimit
	}
	return nextTickSqrtPrice
}

func (s quoteForBaseHelper) ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountQuoteInRemaining math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountQuoteIn := types.CalcAmountQuoteDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, true)

	amountQuoteInAfterFee := amountQuoteInRemaining.Mul(math.LegacyOneDec().Sub(s.feeRate))

	var sqrtPriceNext math.LegacyDec
	if amountQuoteInAfterFee.GTE(amountQuoteIn) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountQuoteInRoundingDown(sqrtPriceCurrent, liquidity, amountQuoteInAfterFee)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountQuoteIn = types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)
	}

	amountBaseOut := types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)

	feeChargeTotal := computeFeeChargePerSwapStepOutGivenIn(hasReachedTarget, amountQuoteIn, amountQuoteInRemaining, s.feeRate)
	return sqrtPriceNext, amountQuoteIn, amountBaseOut, feeChargeTotal
}

func (s quoteForBaseHelper) ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountBaseRemainingOut math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountBaseOut := types.CalcAmountBaseDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, false)

	var sqrtPriceNext math.LegacyDec
	if amountBaseRemainingOut.GTE(amountBaseOut) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountBaseOutRoundingUp(sqrtPriceCurrent, liquidity, amountBaseRemainingOut)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountBaseOut = types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)
	}

	amountQuoteIn := types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)
	feeChargeTotal := computeFeeChargeFromInAmount(amountQuoteIn, getFeeRateOverOneMinusFeeRate(s.feeRate))

	if amountBaseOut.GT(amountBaseRemainingOut) {
		amountBaseOut = amountBaseRemainingOut
	}

	return sqrtPriceNext, amountBaseOut, amountQuoteIn, feeChargeTotal
}

func (s quoteForBaseHelper) NextTickIterator(ctx sdk.Context, storeService store.KVStoreService, poolId uint64, currentTickIndex int64) dbm.Iterator {
	storeAdapter := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
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

func (s quoteForBaseHelper) GetLiquidityDeltaSign(deltaLiquidity math.LegacyDec) math.LegacyDec {
	return deltaLiquidity
}

func (s quoteForBaseHelper) NextTickAfterCrossing(nextTick int64) int64 {
	return nextTick
}

func (s quoteForBaseHelper) ValidateSqrtPrice(sqrtPrice math.LegacyDec, currentSqrtPrice math.LegacyDec) error {
	if sqrtPrice.LT(currentSqrtPrice) || sqrtPrice.GT(types.MaxSqrtPrice) {
		return types.ErrInvalidSqrtPrice
	}
	return nil
}
