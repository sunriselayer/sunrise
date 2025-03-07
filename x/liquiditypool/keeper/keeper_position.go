package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func (k Keeper) initFirstPositionForPool(ctx sdk.Context, pool types.Pool, amountBaseDesired, amountQuoteDesired math.Int) error {
	// Check that the position includes some amount of both base and quote assets
	if !amountBaseDesired.GT(math.ZeroInt()) || !amountQuoteDesired.GT(math.ZeroInt()) {
		return types.ErrInvalidFirstPosition
	}

	// Calculate the spot price and sqrt price
	initialSqrtPrice, err := types.GetSqrtPriceFromQuoteBase(amountQuoteDesired, amountBaseDesired)
	if err != nil {
		return err
	}

	// Calculate the initial tick from the initial spot price
	initialTick, err := types.CalculateSqrtPriceToTick(initialSqrtPrice, pool.TickParams)
	if err != nil {
		return err
	}

	pool.CurrentSqrtPrice = initialSqrtPrice.String()
	pool.CurrentTick = initialTick

	err = k.SetPool(ctx, pool)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to set pool: %d", pool.Id)
	}

	return nil
}

func (k Keeper) resetPool(ctx sdk.Context, pool types.Pool) error {
	pool.CurrentSqrtPrice = math.LegacyZeroDec().String()
	pool.CurrentTick = 0

	err := k.SetPool(ctx, pool)
	if err != nil {
		return errorsmod.Wrapf(err, "failed to set pool: %d", pool.Id)
	}

	return nil
}

func (k Keeper) UpdatePosition(ctx sdk.Context, poolId uint64, owner sdk.AccAddress, lowerTick, upperTick int64, liquidityDelta math.LegacyDec, positionId uint64) (amountBase math.Int, amountQuote math.Int, lowerTickEmpty bool, upperTickEmpty bool, err error) {
	lowerTickIsEmpty, err := k.UpsertTick(ctx, poolId, lowerTick, liquidityDelta, false)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrap(err, "failed to upsert lower tick")
	}

	upperTickIsEmpty, err := k.UpsertTick(ctx, poolId, upperTick, liquidityDelta, true)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrap(err, "failed to upsert upper tick")
	}

	pool, found, err := k.GetPool(ctx, poolId)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(err, "failed to get pool: %d", poolId)
	}
	if !found {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", poolId)
	}

	position, found, err := k.GetPosition(ctx, positionId)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(err, "failed to get position: %d", positionId)
	}
	if !found {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(types.ErrPositionNotFound, "position id: %d", positionId)
	}

	// Update position liquidity
	liquidity, err := math.LegacyNewDecFromStr(position.Liquidity)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}
	liquidity = liquidity.Add(liquidityDelta)
	if liquidity.IsNegative() {
		return math.Int{}, math.Int{}, false, false, types.ErrNegativeLiquidity
	}
	if liquidity.IsZero() {
		err = k.RemovePosition(ctx, position.Id)
		if err != nil {
			return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(err, "failed to remove position: %d", position.Id)
		}
	} else {
		position.Liquidity = liquidity.String()
		err = k.SetPosition(ctx, position)
		if err != nil {
			return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(err, "failed to set position: %d", position.Id)
		}
	}

	actualAmountBase, actualAmountQuote, err := pool.CalcActualAmounts(lowerTick, upperTick, liquidityDelta)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}

	hasPosition, err := k.PoolHasPosition(ctx, poolId)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, errorsmod.Wrapf(err, "failed to check if pool has position: %d", poolId)
	}
	if !hasPosition {
		err = k.resetPool(ctx, pool)
		if err != nil {
			return math.Int{}, math.Int{}, false, false, err
		}
	} else {
		// update liquidity, if position is active
		if pool.IsCurrentTickInRange(lowerTick, upperTick) {
			currentTickLiquidity, err := math.LegacyNewDecFromStr(pool.CurrentTickLiquidity)
			if err != nil {
				return math.Int{}, math.Int{}, false, false, err
			}
			pool.CurrentTickLiquidity = currentTickLiquidity.Add(liquidityDelta).String()
		}

		err = k.SetPool(ctx, pool)
		if err != nil {
			return math.Int{}, math.Int{}, false, false, err
		}
	}

	// update fee accumulator
	if err := k.SetAccumulatorPositionFeeAccumulator(ctx, poolId, lowerTick, upperTick, positionId, liquidityDelta); err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}

	return actualAmountBase.TruncateInt(), actualAmountQuote.TruncateInt(), lowerTickIsEmpty, upperTickIsEmpty, nil
}

func (k Keeper) DecreaseLiquidity(ctx sdk.Context, sender sdk.AccAddress, positionId uint64, liquidity math.LegacyDec) (amountBase math.Int, amountQuote math.Int, err error) {
	// Checks that the element exists
	position, found, err := k.GetPosition(ctx, positionId)
	if err != nil {
		return math.Int{}, math.Int{}, errorsmod.Wrapf(err, "failed to get position: %d", positionId)
	}
	if !found {
		return math.Int{}, math.Int{}, errorsmod.Wrapf(types.ErrPositionNotFound, "id: %d", positionId)
	}

	// Checks if the msg sender is the same as the current owner
	if sender.String() != position.Address {
		return math.Int{}, math.Int{}, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Check if withdrawing negative amount
	if liquidity.IsNegative() {
		return math.Int{}, math.Int{}, types.ErrNegativeTokenAmount
	}

	positionLiquidity, err := math.LegacyNewDecFromStr(position.Liquidity)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}
	if positionLiquidity.LT(liquidity) {
		return math.Int{}, math.Int{}, types.ErrInsufficientLiquidity
	}

	pool, found, err := k.GetPool(ctx, position.PoolId)
	if err != nil {
		return math.Int{}, math.Int{}, errorsmod.Wrapf(err, "failed to get pool: %d", position.PoolId)
	}
	if !found {
		return math.Int{}, math.Int{}, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", position.PoolId)
	}

	// Collect fees
	if _, err := k.collectFees(ctx, sender, positionId); err != nil {
		return math.Int{}, math.Int{}, err
	}

	liquidityDelta := liquidity.Neg()
	amountBase, amountQuote, lowerTickEmpty, upperTickEmpty, err := k.UpdatePosition(ctx, position.PoolId, sender, position.LowerTick, position.UpperTick, liquidityDelta, positionId)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}

	coins := sdk.Coins{sdk.NewCoin(pool.DenomBase, amountBase.Abs())}
	coins = coins.Add(sdk.NewCoin(pool.DenomQuote, amountQuote.Abs()))
	if err := k.bankKeeper.IsSendEnabledCoins(ctx, coins...); err != nil {
		return math.Int{}, math.Int{}, err
	}

	// refund the liquidity to the sender
	err = k.bankKeeper.SendCoins(ctx, pool.GetAddress(), sender, coins)
	if err != nil {
		return math.Int{}, math.Int{}, err
	}

	if lowerTickEmpty {
		k.RemoveTickInfo(ctx, position.PoolId, position.LowerTick)
	}
	if upperTickEmpty {
		k.RemoveTickInfo(ctx, position.PoolId, position.UpperTick)
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventDecreaseLiquidity{
		PositionId:  positionId,
		Address:     sender.String(),
		AmountBase:  amountBase.Abs().String(),
		AmountQuote: amountQuote.Abs().String(),
	}); err != nil {
		return math.Int{}, math.Int{}, err
	}

	return amountBase.Abs(), amountQuote.Abs(), nil
}
