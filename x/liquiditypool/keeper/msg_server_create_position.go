package keeper

import (
	"context"
	"slices"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePosition(ctx context.Context, msg *types.MsgCreatePosition) (*types.MsgCreatePositionResponse, error) {
	senderBytes, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	sender := sdk.AccAddress(senderBytes)
	// end static validation

	// If the sender is one of the authority addresses, skip the sendable token check.
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get params")
	}

	isAllowedSender := slices.Contains(params.AllowedAddresses, sender.String())

	if !isAllowedSender {
		// Validate denom base and denom quote are sendable tokens
		if err := k.bankKeeper.IsSendEnabledCoins(ctx, msg.TokenBase, msg.TokenQuote); err != nil {
			return nil, errorsmod.Wrap(err, "failed to check if send enabled coins")
		}
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	pool, found, err := k.GetPool(ctx, msg.PoolId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get pool")
	}
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", msg.PoolId)
	}

	err = types.CheckTicks(msg.LowerTick, msg.UpperTick)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidTickers, err.Error())
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
		return nil, errorsmod.Wrap(err, "failed to get sqrt price")
	}

	hasPositions := pool.HasPosition(sdkCtx)
	if !hasPositions {
		err := k.initFirstPositionForPool(sdkCtx, pool, amountBaseDesired, amountQuoteDesired)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to init first position for pool")
		}
		// Get updated pool after initialization
		pool, found, err = k.GetPool(ctx, msg.PoolId)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to get pool")
		}
		if !found {
			return nil, errorsmod.Wrapf(types.ErrPoolNotFound, "pool id: %d", msg.PoolId)
		}
	}

	// Calculate the amount of liquidity that will be added to the pool when this position is created.
	currentSqrtPrice, err := math.LegacyNewDecFromStr(pool.CurrentSqrtPrice)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get current sqrt price")
	}
	liquidityDelta := types.GetLiquidityFromAmounts(currentSqrtPrice, sqrtPriceLowerTick, sqrtPriceUpperTick, amountBaseDesired, amountQuoteDesired)
	if liquidityDelta.IsZero() {
		return nil, errorsmod.Wrap(types.ErrZeroLiquidity, `liquidityDelta got 0 value unexpectedly`)
	}

	// Add an empty position in the pool
	var position = types.Position{
		PoolId:    msg.PoolId,
		Address:   msg.Sender,
		LowerTick: msg.LowerTick,
		UpperTick: msg.UpperTick,
		Liquidity: math.LegacyZeroDec().String(),
	}
	positionId, err := k.AppendPosition(ctx, position)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to append position")
	}

	// Initialize / update the position in the pool based on the provided tick range and liquidity delta.
	amountBase, amountQuote, _, _, err := k.UpdatePosition(sdkCtx, position.PoolId, sender, position.LowerTick, position.UpperTick, liquidityDelta, positionId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to update position")
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
	err = k.bankKeeper.SendCoins(ctx, sender, pool.GetAddress(), coins)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins")
	}

	position, _, err = k.GetPosition(ctx, positionId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to get position")
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreatePosition{
		PositionId: positionId,
		Address:    msg.Sender,
		PoolId:     msg.PoolId,
		LowerTick:  msg.LowerTick,
		UpperTick:  msg.UpperTick,
		Liquidity:  position.Liquidity,
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
