package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewPoolAddress(poolId uint64) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("liquiditypool/%d", poolId))
}

func (p Pool) GetAddress() sdk.AccAddress {
	return NewPoolAddress(p.Id)
}

func (p Pool) IsCurrentTickInRange(lowerTick, upperTick int64) bool {
	return p.CurrentTick >= lowerTick && p.CurrentTick < upperTick
}

func (p Pool) CalcActualAmounts(ctx sdk.Context, lowerTick, upperTick int64, liquidityDelta math.LegacyDec) (math.LegacyDec, math.LegacyDec, error) {
	if liquidityDelta.IsZero() {
		return math.LegacyDec{}, math.LegacyDec{}, ErrZeroLiquidity
	}

	sqrtPriceLowerTick, sqrtPriceUpperTick, err := TicksToSqrtPrice(lowerTick, upperTick, p.TickParams)
	if err != nil {
		return math.LegacyDec{}, math.LegacyDec{}, err
	}

	roundUp := liquidityDelta.IsPositive()

	var (
		actualAmountBase  math.LegacyDec
		actualAmountQuote math.LegacyDec
	)

	if p.IsCurrentTickInRange(lowerTick, upperTick) {
		currentSqrtPrice := p.CurrentSqrtPrice
		actualAmountBase = CalcAmountBaseDelta(liquidityDelta, currentSqrtPrice, sqrtPriceUpperTick, roundUp)
		actualAmountQuote = CalcAmountQuoteDelta(liquidityDelta, currentSqrtPrice, sqrtPriceLowerTick, roundUp)
	} else if p.CurrentTick < lowerTick {
		actualAmountQuote = math.LegacyZeroDec()
		actualAmountBase = CalcAmountBaseDelta(liquidityDelta, sqrtPriceLowerTick, sqrtPriceUpperTick, roundUp)
	} else {
		actualAmountBase = math.LegacyZeroDec()
		actualAmountQuote = CalcAmountQuoteDelta(liquidityDelta, sqrtPriceLowerTick, sqrtPriceUpperTick, roundUp)
	}

	if roundUp {
		return actualAmountBase, actualAmountQuote, nil
	}

	return actualAmountBase, actualAmountQuote, nil
}

func (p Pool) HasPosition(ctx sdk.Context) bool {
	if p.CurrentSqrtPrice.IsZero() && p.GetCurrentTick() == 0 {
		return false
	}
	return true
}
