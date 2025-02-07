package keeper

import (
	"context"

	"github.com/sunriselayer/sunrise/x/liquiditypool/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) IncreaseLiquidity(ctx context.Context, msg *types.MsgIncreaseLiquidity) (*types.MsgIncreaseLiquidityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Sender); err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}
	// end static validation

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	position, found := k.GetPosition(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
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

	// Remove full position liquidity
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	amountBaseWithdrawn, amountQuoteWithdrawn, err := k.Keeper.DecreaseLiquidity(sdkCtx, sender, msg.Id, position.Liquidity)
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

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventIncreaseLiquidity{
		OldPositionId: msg.Id,
		NewPositionId: res.Id,
		Address:       msg.Sender,
		AmountBase:    res.AmountBase.String(),
		AmountQuote:   res.AmountQuote.String(),
	}); err != nil {
		return nil, err
	}

	return &types.MsgIncreaseLiquidityResponse{
		PositionId:  res.Id,
		AmountBase:  res.AmountBase,
		AmountQuote: res.AmountQuote,
	}, nil
}
