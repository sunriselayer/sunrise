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

type quoteForBaseStrategy struct {
	sqrtPriceLimit math.LegacyDec
	storeService   store.KVStoreService
	spreadFactor   math.LegacyDec

	oneMinusSpreadFactor math.LegacyDec
	spfOverOneMinusSpf   math.LegacyDec
}

var _ SwapStrategy = (*quoteForBaseStrategy)(nil)

func (s quoteForBaseStrategy) BaseForQuote() bool { return false }

func (s quoteForBaseStrategy) GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec {
	if nextTickSqrtPrice.GT(s.sqrtPriceLimit) {
		return s.sqrtPriceLimit
	}
	return nextTickSqrtPrice
}

func (s quoteForBaseStrategy) ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountQuoteInRemaining math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountQuoteIn := types.CalcAmountQuoteDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, true)

	oneMinusTakerFee := s.getOneMinusSpreadFactor()
	amountQuoteInRemainingLessSpreadReward := amountQuoteInRemaining.Mul(oneMinusTakerFee)

	var sqrtPriceNext math.LegacyDec
	if amountQuoteInRemainingLessSpreadReward.GTE(amountQuoteIn) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountQuoteInRoundingDown(sqrtPriceCurrent, liquidity, amountQuoteInRemainingLessSpreadReward)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountQuoteIn = types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)
	}

	amountBaseOut := types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)

	spreadRewardChargeTotal := computeSpreadRewardChargePerSwapStepOutGivenIn(hasReachedTarget, amountQuoteIn, amountQuoteInRemaining, s.spreadFactor, s.getSpfOverOneMinusSpf)
	return sqrtPriceNext, amountQuoteIn, amountBaseOut, spreadRewardChargeTotal
}

func (s quoteForBaseStrategy) ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountBaseRemainingOut math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
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

	spreadRewardChargeTotal := computeSpreadRewardChargeFromAmountIn(amountQuoteIn, s.getSpfOverOneMinusSpf())

	if amountBaseOut.GT(amountBaseRemainingOut) {
		amountBaseOut = amountBaseRemainingOut
	}

	return sqrtPriceNext, amountBaseOut, amountQuoteIn, spreadRewardChargeTotal
}

func (s quoteForBaseStrategy) getOneMinusSpreadFactor() math.LegacyDec {
	if s.oneMinusSpreadFactor.IsNil() {
		s.oneMinusSpreadFactor = oneDec.Sub(s.spreadFactor)
	}
	return s.oneMinusSpreadFactor
}

func (s quoteForBaseStrategy) getSpfOverOneMinusSpf() math.LegacyDec {
	if s.spfOverOneMinusSpf.IsNil() {
		s.spfOverOneMinusSpf = s.spreadFactor.QuoRoundUp(s.getOneMinusSpreadFactor())
	}
	return s.spfOverOneMinusSpf
}

func (s quoteForBaseStrategy) InitializeNextTickIterator(ctx sdk.Context, poolId uint64, currentTickIndex int64) dbm.Iterator {
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

func (s quoteForBaseStrategy) SetLiquidityDeltaSign(deltaLiquidity math.LegacyDec) math.LegacyDec {
	return deltaLiquidity
}

func (s quoteForBaseStrategy) UpdateTickAfterCrossing(nextTick int64) int64 {
	return nextTick
}

func (s quoteForBaseStrategy) ValidateSqrtPrice(sqrtPrice math.LegacyDec, currentSqrtPrice math.LegacyDec) error {
	if sqrtPrice.LT(currentSqrtPrice) || sqrtPrice.GT(types.MaxSqrtPrice) {
		return types.ErrInvalidSqrtPrice
	}
	return nil
}
