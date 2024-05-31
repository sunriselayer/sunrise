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

func NewPoolFeesAddress(poolId uint64) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("poolfees/%d", poolId))
}

func (p Pool) GetAddress() sdk.AccAddress {
	return NewPoolAddress(p.Id)
}

func (p Pool) GetFeesAddress() sdk.AccAddress {
	return NewPoolFeesAddress(p.Id)
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

	return actualAmountBase, actualAmountQuote, nil
}

func (p Pool) HasPosition(ctx sdk.Context) bool {
	if p.CurrentSqrtPrice.IsZero() && p.GetCurrentTick() == 0 {
		return false
	}
	return true
}

func (p *Pool) UpdateLiquidityIfActivePosition(ctx sdk.Context, lowerTick, upperTick int64, liquidityDelta math.LegacyDec) bool {
	if p.IsCurrentTickInRange(lowerTick, upperTick) {
		p.CurrentTickLiquidity = p.CurrentTickLiquidity.Add(liquidityDelta)
		return true
	}
	return false
}

func (p *Pool) ApplySwap(newLiquidity math.LegacyDec, newCurrentTick int64, newCurrentSqrtPrice math.LegacyDec) error {
	// Check if the new liquidity provided is not negative.
	if newLiquidity.IsNegative() {
		return ErrNegativeLiquidity
	}

	// Check if the new sqrt price provided is not negative.
	if newCurrentSqrtPrice.IsNegative() {
		return ErrNegativeSqrtPrice
	}

	// Check if the new tick provided is within boundaries of the pool's precision factor.
	if newCurrentTick < TICK_MIN || newCurrentTick > TICK_MAX {
		return ErrTickIndexOutOfBoundaries
	}

	p.CurrentTickLiquidity = newLiquidity
	p.CurrentTick = newCurrentTick
	p.CurrentSqrtPrice = newCurrentSqrtPrice

	return nil
}
