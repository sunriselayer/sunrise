package types

import (
	"cosmossdk.io/math"
)

// LiquidityBase = amountBase * (sqrtPriceA * sqrtPriceB) / (sqrtPriceB - sqrtPriceA)
func LiquidityBase(amount math.Int, sqrtPriceA, sqrtPriceB math.LegacyDec) math.LegacyDec {
	if sqrtPriceA.GT(sqrtPriceB) {
		sqrtPriceA, sqrtPriceB = sqrtPriceB, sqrtPriceA
	}

	amounDec := math.LegacyNewDecFromInt(amount)
	product := sqrtPriceA.Mul(sqrtPriceB)
	diff := sqrtPriceB.Sub(sqrtPriceA)
	if diff.IsZero() {
		return math.LegacyZeroDec()
	}

	return amounDec.Mul(product).Quo(diff)
}

// LiquidityQuote = amountQuote / (sqrtPriceB - sqrtPriceA)
func LiquidityQuote(amount math.Int, sqrtPriceA, sqrtPriceB math.LegacyDec) math.LegacyDec {
	if sqrtPriceA.GT(sqrtPriceB) {
		sqrtPriceA, sqrtPriceB = sqrtPriceB, sqrtPriceA
	}

	amountBigDec := math.LegacyNewDecFromInt(amount)
	diff := sqrtPriceB.Sub(sqrtPriceA)
	if diff.IsZero() {
		return math.LegacyZeroDec()
	}

	return amountBigDec.Quo(diff)
}

// CalcAmountBaseDelta = (liquidity * (sqrtPriceB - sqrtPriceA)) / (sqrtPriceB * sqrtPriceA)
func CalcAmountBaseDelta(liq math.LegacyDec, sqrtPriceA, sqrtPriceB math.LegacyDec, roundUp bool) math.LegacyDec {
	if sqrtPriceA.GT(sqrtPriceB) {
		sqrtPriceA, sqrtPriceB = sqrtPriceB, sqrtPriceA
	}
	diff := sqrtPriceB.Sub(sqrtPriceA)
	if roundUp {
		return diff.Mul(liq).Quo(sqrtPriceB).Quo(sqrtPriceA)
	}
	return diff.Mul(liq).Quo(sqrtPriceB).Quo(sqrtPriceA)
}

// CalcAmountQuoteDelta = liq * (sqrtPriceB - sqrtPriceA)
func CalcAmountQuoteDelta(liq math.LegacyDec, sqrtPriceA, sqrtPriceB math.LegacyDec, roundUp bool) math.LegacyDec {
	diff := sqrtPriceB.Sub(sqrtPriceA).AbsMut()
	if roundUp {
		return diff.Mul(liq).Ceil()
	}
	return diff.Mul(liq)
}

// sqrt_next = liq * sqrt_cur / (liq + token_in * sqrt_cur)
func GetNextSqrtPriceFromAmountBaseInRoundingUp(sqrtPriceCurrent, liquidity, amountZeroRemainingIn math.LegacyDec) (sqrtPriceNext math.LegacyDec) {
	if amountZeroRemainingIn.IsZero() {
		return sqrtPriceCurrent
	}

	// Truncate at precision end to make denominator smaller so that the final result is larger.
	product := amountZeroRemainingIn.MulTruncate(sqrtPriceCurrent)
	// denominator = product + liquidity
	denominator := product
	denominator.AddMut(liquidity)
	return liquidity.MulRoundUp(sqrtPriceCurrent).QuoRoundUp(denominator)
}

// sqrt_next = liq * sqrt_cur / (liq - token_out * sqrt_cur)
func GetNextSqrtPriceFromAmountBaseOutRoundingUp(sqrtPriceCurrent, liquidity math.LegacyDec, amountZeroRemainingOut math.LegacyDec) (sqrtPriceNext math.LegacyDec) {
	if amountZeroRemainingOut.IsZero() {
		return sqrtPriceCurrent
	}

	// mul round up to make the final denominator smaller and final result larger
	product := sqrtPriceCurrent.MulRoundUp(amountZeroRemainingOut)
	denominator := liquidity.Sub(product)
	// mul round up numerator to make the final result larger
	// quo round up to make the final result larger
	return liquidity.MulRoundUp(sqrtPriceCurrent).QuoRoundUp(denominator)
}

// sqrt_next = sqrt_cur + token_in / liq
func GetNextSqrtPriceFromAmountQuoteInRoundingDown(sqrtPriceCurrent math.LegacyDec, liquidity math.LegacyDec, amountOneRemainingIn math.LegacyDec) (sqrtPriceNext math.LegacyDec) {
	return amountOneRemainingIn.QuoTruncate(liquidity).Add(sqrtPriceCurrent)
}

// sqrt_next = sqrt_cur - token_out / liq
func GetNextSqrtPriceFromAmountQuoteOutRoundingDown(sqrtPriceCurrent math.LegacyDec, liquidity math.LegacyDec, amountOneRemainingOut math.LegacyDec) (sqrtPriceNext math.LegacyDec) {
	return sqrtPriceCurrent.Sub(amountOneRemainingOut.QuoRoundUp(liquidity))
}

func GetLiquidityFromAmounts(sqrtPrice math.LegacyDec, sqrtPriceA, sqrtPriceB math.LegacyDec, amountBase, amountQuote math.Int) (liquidity math.LegacyDec) {
	if sqrtPriceA.GT(sqrtPriceB) {
		sqrtPriceA, sqrtPriceB = sqrtPriceB, sqrtPriceA
	}

	if sqrtPrice.LTE(sqrtPriceA) {
		liquidity = LiquidityBase(amountBase, sqrtPriceA, sqrtPriceB)
	} else if sqrtPrice.LT(sqrtPriceB) {
		liquidityBase := LiquidityBase(amountBase, sqrtPrice, sqrtPriceB)
		liquidityQuote := LiquidityQuote(amountQuote, sqrtPrice, sqrtPriceA)
		liquidity = math.LegacyMinDec(liquidityBase, liquidityQuote)
	} else {
		liquidity = LiquidityQuote(amountQuote, sqrtPriceB, sqrtPriceA)
	}

	return liquidity
}

func SquareRoundUp(sqrtPrice math.LegacyDec) math.LegacyDec {
	return sqrtPrice.MulRoundUp(sqrtPrice)
}

func SquareTruncate(sqrtPrice math.LegacyDec) math.LegacyDec {
	return sqrtPrice.MulTruncate(sqrtPrice)
}
