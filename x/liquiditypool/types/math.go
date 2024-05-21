package types

import (
	"fmt"

	"cosmossdk.io/math"
)

var powPrecision, _ = math.LegacyNewDecFromStr("0.00000001")

var (
	oneHalf math.LegacyDec = math.LegacyMustNewDecFromStr("0.5")
	one     math.LegacyDec = math.LegacyOneDec()
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

func Pow(base math.LegacyDec, exp math.LegacyDec) math.LegacyDec {
	if !base.IsPositive() {
		panic(fmt.Errorf("base must be greater than 0"))
	}

	// We will use an approximation algorithm to compute the power.
	// Since computing an integer power is easy, we split up the exponent into
	// an integer component and a fractional component.
	integer := exp.TruncateDec()
	fractional := exp.Sub(integer)

	integerPow := base.Power(uint64(integer.TruncateInt64()))

	if fractional.IsZero() {
		return integerPow
	}

	fractionalPow := PowApprox(base, fractional, powPrecision)

	return integerPow.Mul(fractionalPow)
}

// Contract: 0 < base <= 2
// 0 <= exp < 1.
func PowApprox(originalBase math.LegacyDec, exp math.LegacyDec, precision math.LegacyDec) math.LegacyDec {
	if !originalBase.IsPositive() {
		panic(fmt.Errorf("base must be greater than 0"))
	}

	if exp.IsZero() {
		return math.LegacyOneDec()
	}

	// Common case optimization
	// Optimize for it being equal to one-half
	if exp.Equal(oneHalf) {
		output, err := originalBase.ApproxSqrt()
		if err != nil {
			panic(err)
		}
		return output
	}

	base := originalBase.Clone()
	x, xneg := AbsDifferenceWithSign(base, one)
	term := math.LegacyOneDec()
	sum := math.LegacyOneDec()
	negative := false

	a := exp.Clone()
	bigK := math.LegacyNewDec(0)
	for i := int64(1); term.GTE(precision); i++ {
		// At each iteration, we need two values, i and i-1.
		// To avoid expensive big.Int allocation, we reuse bigK variable.
		// On this line, bigK == i-1.
		c, cneg := AbsDifferenceWithSign(a, bigK)
		// On this line, bigK == i.
		bigK.SetInt64(i)
		term.MulMut(c).MulMut(x).QuoMut(bigK)

		// a is mutated on absDifferenceWithSign, reset
		a.Set(exp)

		if term.IsZero() {
			break
		}
		if xneg {
			negative = !negative
		}

		if cneg {
			negative = !negative
		}

		if negative {
			sum.SubMut(term)
		} else {
			sum.AddMut(term)
		}
	}
	return sum
}

// AbsDifferenceWithSign returns | a - b |, (a - b).sign()
// a is mutated and returned.
func AbsDifferenceWithSign(a, b math.LegacyDec) (math.LegacyDec, bool) {
	if a.GTE(b) {
		return a.SubMut(b), false
	} else {
		return a.NegMut().AddMut(b), true
	}
}
