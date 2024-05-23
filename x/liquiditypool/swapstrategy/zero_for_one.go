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

type zeroForOneStrategy struct {
	sqrtPriceLimit math.LegacyDec
	storeService   store.KVStoreService
	spreadFactor   math.LegacyDec

	oneMinusSpreadFactor math.LegacyDec
	spfOverOneMinusSpf   math.LegacyDec
}

var _ SwapStrategy = (*zeroForOneStrategy)(nil)

func (s zeroForOneStrategy) ZeroForOne() bool { return true }

func (s zeroForOneStrategy) GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec {
	if nextTickSqrtPrice.LT(s.sqrtPriceLimit) {
		return s.sqrtPriceLimit
	}
	return nextTickSqrtPrice
}

func (s zeroForOneStrategy) ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountZeroInRemaining math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountZeroIn := types.CalcAmountBaseDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, true) // N.B.: if this is false, causes infinite loop

	oneMinusTakerFee := s.getOneMinusSpreadFactor()
	amountZeroInRemainingLessSpreadReward := amountZeroInRemaining.Mul(oneMinusTakerFee)

	var sqrtPriceNext math.LegacyDec
	if amountZeroInRemainingLessSpreadReward.GTE(amountZeroIn) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountBaseInRoundingUp(sqrtPriceCurrent, liquidity, amountZeroInRemainingLessSpreadReward)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountZeroIn = types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true) // N.B.: if this is false, causes infinite loop
	}

	amountOneOut := types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)

	spreadRewardChargeTotal := computeSpreadRewardChargePerSwapStepOutGivenIn(hasReachedTarget, amountZeroIn, amountZeroInRemaining, s.spreadFactor, s.getSpfOverOneMinusSpf)

	return sqrtPriceNext, amountZeroIn, amountOneOut, spreadRewardChargeTotal
}

func (s zeroForOneStrategy) ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountOneRemainingOut math.LegacyDec) (math.LegacyDec, math.LegacyDec, math.LegacyDec, math.LegacyDec) {
	amountOneOut := types.CalcAmountQuoteDelta(liquidity, sqrtPriceTarget, sqrtPriceCurrent, false)

	var sqrtPriceNext math.LegacyDec
	if amountOneRemainingOut.GTE(amountOneOut) {
		sqrtPriceNext = sqrtPriceTarget
	} else {
		sqrtPriceNext = types.GetNextSqrtPriceFromAmountQuoteOutRoundingDown(sqrtPriceCurrent, liquidity, amountOneRemainingOut)
	}

	hasReachedTarget := sqrtPriceTarget.Equal(sqrtPriceNext)

	if !hasReachedTarget {
		amountOneOut = types.CalcAmountQuoteDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, false)
	}

	amountZeroIn := types.CalcAmountBaseDelta(liquidity, sqrtPriceNext, sqrtPriceCurrent, true)

	spreadRewardChargeTotal := computeSpreadRewardChargeFromAmountIn(amountZeroIn, s.getSpfOverOneMinusSpf())

	if amountOneOut.GT(amountOneRemainingOut) {
		amountOneOut = amountOneRemainingOut
	}

	return sqrtPriceNext, amountOneOut, amountZeroIn, spreadRewardChargeTotal
}

func (s zeroForOneStrategy) getOneMinusSpreadFactor() math.LegacyDec {
	if s.oneMinusSpreadFactor.IsNil() {
		s.oneMinusSpreadFactor = oneDec.Sub(s.spreadFactor)
	}
	return s.oneMinusSpreadFactor
}

func (s zeroForOneStrategy) getSpfOverOneMinusSpf() math.LegacyDec {
	if s.spfOverOneMinusSpf.IsNil() {
		s.spfOverOneMinusSpf = s.spreadFactor.QuoRoundUp(s.getOneMinusSpreadFactor())
	}
	return s.spfOverOneMinusSpf
}

func (s zeroForOneStrategy) InitializeNextTickIterator(ctx sdk.Context, poolId uint64, currentTickIndex int64) dbm.Iterator {
	storeAdapter := runtime.KVStoreAdapter(s.storeService.OpenKVStore(ctx))
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

func (s zeroForOneStrategy) SetLiquidityDeltaSign(deltaLiquidity math.LegacyDec) math.LegacyDec {
	return deltaLiquidity.Neg()
}

func (s zeroForOneStrategy) UpdateTickAfterCrossing(nextTick int64) int64 {
	return nextTick - 1
}

func (s zeroForOneStrategy) ValidateSqrtPrice(sqrtPrice math.LegacyDec, currentSqrtPrice math.LegacyDec) error {
	if sqrtPrice.GT(currentSqrtPrice) || sqrtPrice.LT(types.MinSqrtPrice) {
		return types.ErrInvalidSqrtPrice
	}
	return nil
}
