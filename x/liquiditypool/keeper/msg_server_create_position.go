package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePosition(ctx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	// end static validation
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", msg.PoolId)
	}

	err = types.CheckTicks(msg.LowerTick, msg.UpperTick)
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

	hasPositions := pool.HasPosition(sdkCtx)
	if !hasPositions {
		err := k.initFirstPositionForPool(sdkCtx, pool, amountBaseDesired, amountQuoteDesired)
		if err != nil {
			return nil, err
		}
		// Get updated pool after initialization
		pool, _ = k.GetPool(ctx, msg.PoolId)
	}

	// Calculate the amount of liquidity that will be added to the pool when this position is created.
	liquidityDelta := types.GetLiquidityFromAmounts(pool.CurrentSqrtPrice, sqrtPriceLowerTick, sqrtPriceUpperTick, amountBaseDesired, amountQuoteDesired)
	if liquidityDelta.IsZero() {
		return nil, errorsmod.Wrap(types.ErrZeroLiquidity, `liquidityDelta got 0 value unexpectedly`)
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
	amountBase, amountQuote, _, _, err := k.UpdatePosition(sdkCtx, position.PoolId, sender, position.LowerTick, position.UpperTick, liquidityDelta, positionId)
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

	position, _ = k.GetPosition(ctx, positionId)

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreatePosition{
		PositionId: positionId,
		Address:    msg.Sender,
		PoolId:     msg.PoolId,
		LowerTick:  msg.LowerTick,
		UpperTick:  msg.UpperTick,
		Liquidity:  position.Liquidity.String(),
	}); err != nil {
		return nil, err
	}

	return &types.MsgCreatePositionResponse{
		Id:          positionId,
		AmountBase:  amountBase,
		AmountQuote: amountQuote,
		Liquidity:   position.Liquidity,
	}, nil
}
