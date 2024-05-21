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

	if tickUpper < types.TICK_MAX {
		return errors.New("tickUpper is out of range")
	}
	return nil
}

func (k msgServer) CreatePosition(goCtx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", msg.PoolId)
	}

	err := checkTicks(msg.LowerTick, msg.UpperTick)
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

	// Transform the provided ticks to sqrtPrices.
	sqrtPriceLowerTick, sqrtPriceUpperTick, err := types.TicksToSqrtPrice(msg.LowerTick, msg.UpperTick, pool.TickParams)
	if err != nil {
		return nil, err
	}

	// If this is the first position created in this pool, ensure that the position includes both assets.
	hasPositions := pool.HasPosition(ctx)
	if !hasPositions {
		err := k.initFirstPositionForPool(ctx, pool, amountBaseDesired, amountQuoteDesired)
		if err != nil {
			return nil, err
		}
	}

	// Calculate the amount of liquidity that will be added to the pool when this position is created.
	liquidityDelta := types.GetLiquidityFromAmounts(pool.CurrentSqrtPrice, sqrtPriceLowerTick, sqrtPriceUpperTick, amountBaseDesired, amountQuoteDesired)
	if liquidityDelta.IsZero() {
		return nil, fmt.Errorf(`liquidityDelta got 0 value unexpectedly`)
	}

	// Add a position in the pool
	var position = types.Position{
		PoolId:    msg.PoolId,
		Address:   msg.Sender,
		LowerTick: msg.LowerTick,
		UpperTick: msg.UpperTick,
		Liquidity: liquidityDelta,
	}

	id := k.AppendPosition(ctx, position)

	// calculate the actual amounts of tokens 0 and 1 that were added or removed from the pool.
	amountBase, amountQuote, err := pool.CalcActualAmounts(ctx, msg.LowerTick, msg.UpperTick, liquidityDelta)
	if err != nil {
		return nil, err
	}

	// Check if the actual amounts are greater than or equal to minimum amounts
	if amountBase.LT(msg.MinAmountBase.ToLegacyDec()) {
		return nil, types.ErrInsufficientAmountPut
	}
	if amountQuote.LT(msg.MinAmountQuote.ToLegacyDec()) {
		return nil, types.ErrInsufficientAmountPut
	}

	// Transfer amounts to the pool
	coins := sdk.Coins{}
	amountBaseInt := amountBase.TruncateInt()
	amountQuoteInt := amountQuote.TruncateInt()
	coins = coins.Add(sdk.NewCoin(msg.TokenBase.Denom, amountBaseInt))
	coins = coins.Add(sdk.NewCoin(msg.TokenQuote.Denom, amountQuoteInt))
	err = k.bankKeeper.SendCoins(ctx, sender, pool.GetAddress(), coins)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreatePositionResponse{
		Id:          id,
		AmountBase:  amountBaseInt,
		AmountQuote: amountQuoteInt,
		Liquidity:   position.Liquidity,
	}, nil
}

func (k msgServer) IncreaseLiquidity(goCtx context.Context, msg *types.MsgIncreaseLiquidity) (*types.MsgIncreaseLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var position = types.Position{
		Address: msg.Sender,
		Id:      msg.Id,
	}

	// Checks that the element exists
	val, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != val.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPosition(ctx, position)

	return &types.MsgIncreaseLiquidityResponse{}, nil
}

func (k msgServer) DecreaseLiquidity(goCtx context.Context, msg *types.MsgDecreaseLiquidity) (*types.MsgDecreaseLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	position, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg sender is the same as the current owner
	if msg.Sender != position.Address {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePosition(ctx, msg.Id)

	return &types.MsgDecreaseLiquidityResponse{}, nil
}

func (k Keeper) initFirstPositionForPool(ctx sdk.Context, pool types.Pool, amountBaseDesired, amountQuoteDesired math.Int) error {
	// Check that the position includes some amount of both base and quote assets
	if !amountBaseDesired.GT(math.ZeroInt()) || !amountQuoteDesired.GT(math.ZeroInt()) {
		return types.ErrInvalidFirstPosition
	}

	// Calculate the spot price and sqrt price
	initialSpotPrice := amountQuoteDesired.ToLegacyDec().Quo(amountBaseDesired.ToLegacyDec())
	initialCurSqrtPrice, err := initialSpotPrice.ApproxSqrt()
	if err != nil {
		return err
	}

	// Calculate the initial tick from the initial spot price
	initialTick, err := types.CalculateSqrtPriceToTick(initialCurSqrtPrice, pool.TickParams)
	if err != nil {
		return err
	}

	pool.CurrentSqrtPrice = initialCurSqrtPrice
	pool.CurrentTick = initialTick
	k.SetPool(ctx, pool)

	return nil
}
