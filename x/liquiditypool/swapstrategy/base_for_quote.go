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

type baseForQuoteStrategy struct {
	sqrtPriceLimit math.LegacyDec
	storeService   store.KVStoreService
	spreadFactor   math.LegacyDec

	oneMinusSpreadFactor math.LegacyDec
	spfOverOneMinusSpf   math.LegacyDec
}

var _ SwapStrategy = (*baseForQuoteStrategy)(nil)

func (s baseForQuoteStrategy) BaseForQuote() bool { return true }

func (s baseForQuoteStrategy) GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec {
	if nextTickSqrtPrice.LT(s.sqrtPriceLimit) {
		return s.sqrtPriceLimit
	}
	return nextTickSqrtPrice
}

func (s baseForQuoteStrategy) ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountBaseInRemaining math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountBaseIn := types.CalcAmountBaseDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, true)
	oneMinusTakerFee := s.getOneMinusSpreadFactor()
	amountBaseInRemainingLessSpreadReward := amountBaseInRemaining.Mul(oneMinusTakerFee)

	var sqrtPriceNext math.LegacyDec
	if amountBaseInRemainingLessSpreadReward.GTE(amountBaseIn) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountBaseInRoundingUp(sqrtPriceCurrent, liquidity, amountBaseInRemainingLessSpreadReward)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)
	if !hasReachedTarget {
		amountBaseIn = types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)
	}

	amountQuoteOut := types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)

	spreadRewardChargeTotal := computeSpreadRewardChargePerSwapStepOutGivenIn(hasReachedTarget, amountBaseIn, amountBaseInRemaining, s.spreadFactor, s.getSpfOverOneMinusSpf)
	return sqrtPriceNext, amountBaseIn, amountQuoteOut, spreadRewardChargeTotal
}

func (s baseForQuoteStrategy) ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountQuoteRemainingOut math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
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

	spreadRewardChargeTotal := computeSpreadRewardChargeFromAmountIn(amountBaseIn, s.getSpfOverOneMinusSpf())

	if amountQuoteOut.GT(amountQuoteRemainingOut) {
		amountQuoteOut = amountQuoteRemainingOut
	}

	return sqrtPriceNext, amountQuoteOut, amountBaseIn, spreadRewardChargeTotal
}

func (s baseForQuoteStrategy) getOneMinusSpreadFactor() math.LegacyDec {
	if s.oneMinusSpreadFactor.IsNil() {
		s.oneMinusSpreadFactor = oneDec.Sub(s.spreadFactor)
	}
	return s.oneMinusSpreadFactor
}

func (s baseForQuoteStrategy) getSpfOverOneMinusSpf() math.LegacyDec {
	if s.spfOverOneMinusSpf.IsNil() {
		s.spfOverOneMinusSpf = s.spreadFactor.QuoRoundUp(s.getOneMinusSpreadFactor())
	}
	return s.spfOverOneMinusSpf
}

func (s baseForQuoteStrategy) InitializeNextTickIterator(ctx sdk.Context, poolId uint64, currentTickIndex int64) dbm.Iterator {
	storeAdapter := runtime.KVStoreAdapter(s.storeService.OpenKVStore(ctx))
	// store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.TickInfoKey))
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

func (s baseForQuoteStrategy) SetLiquidityDeltaSign(deltaLiquidity math.LegacyDec) math.LegacyDec {
	return deltaLiquidity.Neg()
}

func (s baseForQuoteStrategy) UpdateTickAfterCrossing(nextTick int64) int64 {
	return nextTick - 1
}

func (s baseForQuoteStrategy) ValidateSqrtPrice(sqrtPrice math.LegacyDec, currentSqrtPrice math.LegacyDec) error {
	if sqrtPrice.GT(currentSqrtPrice) || sqrtPrice.LT(types.MinSqrtPrice) {
		return types.ErrInvalidSqrtPrice
	}
	return nil
}
