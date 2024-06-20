package keeper

import (
	"context"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

func checkTicks(tickLower, tickUpper int64) error {
	if tickLower >= tickUpper {
		return errors.New("tickLower should be less than tickUpper")
	}
	if tickLower < types.TICK_MIN {
		return errors.New("tickLower is out of range")
	}

	if tickUpper > types.TICK_MAX {
		return errors.New("tickUpper is out of range")
	}
	return nil
}

func (k msgServer) CreatePosition(goCtx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", msg.PoolId)
	}

	err = checkTicks(msg.LowerTick, msg.UpperTick)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidTickers, err.Error())
	}

	if pool.DenomBase != msg.TokenBase.Denom {
		return nil, errorsmod.Wrapf(types.ErrInvalidBaseDenom, "expected %s, provided %s", pool.DenomBase, msg.TokenBase.Denom)
	}

	if pool.DenomQuote != msg.TokenQuote.Denom {
		return nil, errorsmod.Wrapf(types.ErrInvalidQuoteDenom, "expected %s, provided %s", pool.DenomQuote, msg.TokenQuote.Denom)
	}

	if msg.TokenBase.Amount.IsZero() && msg.TokenQuote.Amount.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrInvalidTokenAmounts, "base amount %s, quote amount %s", msg.TokenBase.String(), msg.TokenQuote.String())
	}

	if msg.TokenBase.Amount.IsNegative() || msg.TokenQuote.Amount.IsNegative() {
		return nil, types.ErrNegativeTokenAmount
	}

	amountBaseDesired := msg.TokenBase.Amount
	amountQuoteDesired := msg.TokenQuote.Amount

	sqrtPriceLowerTick, sqrtPriceUpperTick, err := types.TicksToSqrtPrice(msg.LowerTick, msg.UpperTick, pool.TickParams)
	if err != nil {
		return nil, err
	}

	hasPositions := pool.HasPosition(ctx)
	if !hasPositions {
		err := k.initFirstPositionForPool(ctx, pool, amountBaseDesired, amountQuoteDesired)
		if err != nil {
			return nil, err
		}
		// Get updated pool after initialization
		pool, _ = k.GetPool(ctx, msg.PoolId)
	}

	// Calculate the amount of liquidity that will be added to the pool when this position is created.
	liquidityDelta := types.GetLiquidityFromAmounts(pool.CurrentSqrtPrice, sqrtPriceLowerTick, sqrtPriceUpperTick, amountBaseDesired, amountQuoteDesired)
	if liquidityDelta.IsZero() {
		return nil, fmt.Errorf(`liquidityDelta got 0 value unexpectedly`)
	}

	// Add an empty position in the pool
	var position = types.Position{
		PoolId:    msg.PoolId,
		Address:   msg.Sender,
		LowerTick: msg.LowerTick,
		UpperTick: msg.UpperTick,
		Liquidity: math.LegacyZeroDec(),
	}
	positionId := k.AppendPosition(ctx, position)

	// Initialize / update the position in the pool based on the provided tick range and liquidity delta.
	amountBase, amountQuote, _, _, err := k.UpdatePosition(ctx, position.PoolId, sender, position.LowerTick, position.UpperTick, liquidityDelta, positionId)
	if err != nil {
		return nil, err
	}

	// Check if the actual amounts are greater than or equal to minimum amounts
	if amountBase.LT(msg.MinAmountBase) {
		return nil, errorsmod.Wrapf(types.ErrInsufficientAmountPut, "min_base_amount; expected %s, got %s", amountBase, msg.MinAmountBase)
	}
	if amountQuote.LT(msg.MinAmountQuote) {
		return nil, errorsmod.Wrapf(types.ErrInsufficientAmountPut, "min_quote_amount; expected %s, got %s", amountQuote, msg.MinAmountQuote)
	}

	// Transfer amounts to the pool
	coins := sdk.Coins{sdk.NewCoin(msg.TokenBase.Denom, amountBase)}
	coins = coins.Add(sdk.NewCoin(msg.TokenQuote.Denom, amountQuote))
	if err := k.bankKeeper.IsSendEnabledCoins(ctx, coins...); err != nil {
		return nil, err
	}
	err = k.bankKeeper.SendCoins(ctx, sender, pool.GetAddress(), coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreatePositionResponse{
		Id:          positionId,
		AmountBase:  amountBase,
		AmountQuote: amountQuote,
		Liquidity:   position.Liquidity,
	}, nil
}

func (k msgServer) IncreaseLiquidity(goCtx context.Context, msg *types.MsgIncreaseLiquidity) (*types.MsgIncreaseLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	position, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	if msg.Sender != position.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if msg.AmountBase.IsNegative() || msg.AmountQuote.IsNegative() {
		return nil, types.ErrNegativeTokenAmount
	}

	if msg.AmountBase.IsZero() && msg.AmountQuote.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrInvalidTokenAmounts, "base amount %s, quote amount %s", msg.AmountBase.String(), msg.AmountQuote.String())
	}

	k.SetPosition(ctx, position)

	// Remove full position liquidity
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	amountBaseWithdrawn, amountQuoteWithdrawn, err := k.Keeper.DecreaseLiquidity(ctx, sender, msg.Id, position.Liquidity)
	if err != nil {
		return nil, err
	}

	pool, found := k.GetPool(ctx, position.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	// Create a new position with combined liquidity
	amountBaseDesired := amountBaseWithdrawn.Add(msg.AmountBase)
	amountQuoteDesired := amountQuoteWithdrawn.Add(msg.AmountQuote)
	minAmountBase := amountBaseWithdrawn.Add(msg.MinAmountBase)
	minAmountQuote := amountQuoteWithdrawn.Add(msg.MinAmountQuote)

	res, err := k.CreatePosition(ctx, &types.MsgCreatePosition{
		Sender:         msg.Sender,
		PoolId:         position.PoolId,
		LowerTick:      position.LowerTick,
		UpperTick:      position.UpperTick,
		TokenBase:      sdk.NewCoin(pool.DenomBase, amountBaseDesired),
		TokenQuote:     sdk.NewCoin(pool.DenomQuote, amountQuoteDesired),
		MinAmountBase:  minAmountBase,
		MinAmountQuote: minAmountQuote,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgIncreaseLiquidityResponse{
		AmountBase:  res.AmountBase,
		AmountQuote: res.AmountQuote,
	}, nil
}

func (k Keeper) DecreaseLiquidity(ctx sdk.Context, sender sdk.AccAddress, positionId uint64, liquidity math.LegacyDec) (amountBase math.Int, amountQuote math.Int, err error) {
	// Checks that the element exists
	position, found := k.GetPosition(ctx, positionId)
	if !found {
		return math.Int{}, math.Int{}, errorsmod.Wrap(types.ErrPositionNotFound, fmt.Sprintf("id: %d", positionId))
	}

	// Checks if the msg sender is the same as the current owner
	if sender.String() != position.Address {
		return math.Int{}, math.Int{}, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Check if withdrawing negative amount
	if liquidity.IsNegative() {
		return math.Int{}, math.Int{}, types.ErrNegativeTokenAmount
	}

	if position.Liquidity.LT(liquidity) {
		return math.Int{}, math.Int{}, types.ErrInsufficientLiquidity
	}

	pool, found := k.GetPool(ctx, position.PoolId)
	if !found {
		return math.Int{}, math.Int{}, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", position.PoolId)
	}

	liquiditDelta := liquidity.Neg()
	amountBase, amountQuote, lowerTickEmpty, upperTickEmpty, err := k.UpdatePosition(ctx, position.PoolId, sender, position.LowerTick, position.UpperTick, liquiditDelta, positionId)
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

	if liquidity.Equal(position.Liquidity) {
		// Collect fees
		if _, err := k.collectFees(ctx, sender, positionId); err != nil {
			return math.Int{}, math.Int{}, err
		}
		k.RemovePosition(ctx, position.Id)

		if !k.PoolHasPosition(ctx, position.PoolId) {
			k.resetPool(ctx, pool)
		}
	}

	if lowerTickEmpty {
		k.RemoveTickInfo(ctx, position.PoolId, position.LowerTick)
	}
	if upperTickEmpty {
		k.RemoveTickInfo(ctx, position.PoolId, position.UpperTick)
	}
	return amountBase.Abs(), amountQuote.Abs(), nil
}

func (k msgServer) DecreaseLiquidity(goCtx context.Context, msg *types.MsgDecreaseLiquidity) (*types.MsgDecreaseLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	amountBase, amountQuote, err := k.Keeper.DecreaseLiquidity(ctx, sender, msg.Id, msg.Liquidity)
	if err != nil {
		return nil, err
	}

	return &types.MsgDecreaseLiquidityResponse{
		AmountBase:  amountBase,
		AmountQuote: amountQuote,
	}, nil
}

func (k Keeper) initFirstPositionForPool(ctx sdk.Context, pool types.Pool, amountBaseDesired, amountQuoteDesired math.Int) error {
	// Check that the position includes some amount of both base and quote assets
	if !amountBaseDesired.GT(math.ZeroInt()) || !amountQuoteDesired.GT(math.ZeroInt()) {
		return types.ErrInvalidFirstPosition
	}

	// Calculate the spot price and sqrt price
	initialSpotPrice := amountQuoteDesired.ToLegacyDec().Quo(amountBaseDesired.ToLegacyDec())
	initialCurrentSqrtPrice, err := initialSpotPrice.ApproxSqrt()
	if err != nil {
		return err
	}

	// Calculate the initial tick from the initial spot price
	initialTick, err := types.CalculateSqrtPriceToTick(initialCurrentSqrtPrice, pool.TickParams)
	if err != nil {
		return err
	}

	pool.CurrentSqrtPrice = initialCurrentSqrtPrice
	pool.CurrentTick = initialTick
	k.SetPool(ctx, pool)

	return nil
}

func (k Keeper) resetPool(ctx sdk.Context, pool types.Pool) {
	pool.CurrentSqrtPrice = math.LegacyZeroDec()
	pool.CurrentTick = 0

	k.SetPool(ctx, pool)
}

func (k Keeper) UpdatePosition(ctx sdk.Context, poolId uint64, owner sdk.AccAddress, lowerTick, upperTick int64, liquidityDelta math.LegacyDec, positionId uint64) (amountBase math.Int, amountQuote math.Int, lowerTickEmpty bool, upperTickEmpty bool, err error) {
	lowerTickIsEmpty, err := k.initOrUpdateTick(ctx, poolId, lowerTick, liquidityDelta, false)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}

	upperTickIsEmpty, err := k.initOrUpdateTick(ctx, poolId, upperTick, liquidityDelta, true)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return math.Int{}, math.Int{}, false, false, types.ErrPoolNotFound
	}

	position, found := k.GetPosition(ctx, positionId)
	if !found {
		return math.Int{}, math.Int{}, false, false, types.ErrPositionNotFound
	}

	// Update position liquidity
	position.Liquidity = position.Liquidity.Add(liquidityDelta)
	if position.Liquidity.IsNegative() {
		return math.Int{}, math.Int{}, false, false, types.ErrNegativeLiquidity
	}
	k.SetPosition(ctx, position)

	actualAmountBase, actualAmountQuote, err := pool.CalcActualAmounts(lowerTick, upperTick, liquidityDelta)
	if err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}

	pool.UpdateLiquidityIfActivePosition(ctx, lowerTick, upperTick, liquidityDelta)

	k.SetPool(ctx, pool)

	// update fee accumulator
	if err := k.SetAccumulatorPositionFeeAccumulator(ctx, poolId, lowerTick, upperTick, positionId, liquidityDelta); err != nil {
		return math.Int{}, math.Int{}, false, false, err
	}

	return actualAmountBase.TruncateInt(), actualAmountQuote.TruncateInt(), lowerTickIsEmpty, upperTickIsEmpty, nil
}
