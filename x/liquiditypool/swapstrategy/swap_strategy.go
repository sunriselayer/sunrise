package swapstrategy

import (
	dbm "github.com/cometbft/cometbft-db"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	"cosmossdk.io/core/store"
)

type SwapStrategy interface {
	GetSqrtTargetPrice(nextTickSqrtPrice math.LegacyDec) math.LegacyDec
	ComputeSwapWithinBucketOutGivenIn(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountRemainingIn math.LegacyDec) (sqrtPriceNext math.LegacyDec, amountInConsumed, amountOutComputed, spreadRewardChargeTotal math.LegacyDec)
	ComputeSwapWithinBucketInGivenOut(sqrtPriceCurrent, sqrtPriceTarget math.LegacyDec, liquidity, amountRemainingOut math.LegacyDec) (sqrtPriceNext math.LegacyDec, amountOutConsumed, amountInComputed, spreadRewardChargeTotal math.LegacyDec)
	InitializeNextTickIterator(ctx sdk.Context, poolId uint64, tickIndex int64) dbm.Iterator
	SetLiquidityDeltaSign(liquidityDelta math.LegacyDec) math.LegacyDec
	UpdateTickAfterCrossing(nextTick int64) (updatedNextTick int64)
	ValidateSqrtPrice(sqrtPriceLimit math.LegacyDec, currentSqrtPrice math.LegacyDec) error
	BaseForQuote() bool
}

var (
	oneDec = math.LegacyOneDec()
)

func New(baseForQuote bool, sqrtPriceLimit math.LegacyDec, storeService store.KVStoreService, spreadFactor math.LegacyDec) SwapStrategy {
	if baseForQuote {
		return &baseForQuoteStrategy{sqrtPriceLimit: sqrtPriceLimit, storeService: storeService, spreadFactor: spreadFactor}
	}
	return &quoteForBaseStrategy{sqrtPriceLimit: sqrtPriceLimit, storeService: storeService, spreadFactor: spreadFactor}
}

func GetPriceLimit(baseForQuote bool) math.LegacyDec {
	if baseForQuote {
		return types.MinSpotPrice
	}
	return types.MaxSpotPrice
}

func GetSqrtPriceLimit(priceLimit math.LegacyDec, baseForQuote bool) (math.LegacyDec, error) {
	if priceLimit.IsZero() {
		if baseForQuote {
			return types.MinSqrtPrice, nil
		}
		return types.MaxSqrtPrice, nil
	}

	if priceLimit.LT(types.MinSpotPrice) || priceLimit.GT(types.MaxSpotPrice) {
		return math.LegacyDec{}, types.ErrPriceOutOfBound
	}

	sqrtPriceLimit, err := priceLimit.ApproxSqrt()
	if err != nil {
		return math.LegacyDec{}, err
	}

	return sqrtPriceLimit, nil
}
