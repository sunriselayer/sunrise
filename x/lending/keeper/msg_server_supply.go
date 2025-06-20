package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func (k msgServer) Supply(ctx context.Context, msg *types.MsgSupply) (*types.MsgSupplyResponse, error) {
	// Validate sender address
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Validate amount
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "supply amount must be positive")
	}

	// Check if market exists, create if not
	denom := msg.Amount.Denom
	market, err := k.Markets.Get(ctx, denom)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			// Create new market
			market = types.Market{
				Denom:             denom,
				TotalSupplied:     math.ZeroInt(),
				TotalBorrowed:     math.ZeroInt(),
				GlobalRewardIndex: math.LegacyOneDec(),
				RiseDenom:         fmt.Sprintf("rise%s", denom),
			}
		} else {
			return nil, err
		}
	}

	// Transfer tokens from sender to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.AccAddress(sender),
		types.ModuleName,
		sdk.NewCoins(msg.Amount),
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to transfer tokens to module")
	}

	// Update market total supplied
	market.TotalSupplied = market.TotalSupplied.Add(msg.Amount.Amount)
	if err := k.Markets.Set(ctx, denom, market); err != nil {
		return nil, err
	}

	// Get or create user position
	position, err := k.UserPositions.Get(ctx, collections.Join(msg.Sender, denom))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			// Create new position
			position = types.UserPosition{
				UserAddress:      msg.Sender,
				Denom:            denom,
				Amount:           math.ZeroInt(),
				LastRewardIndex:  market.GlobalRewardIndex,
			}
		} else {
			return nil, err
		}
	}

	// Calculate rise tokens to mint (1:1 for now, will add interest calculation later)
	riseAmount := msg.Amount.Amount

	// Update user position
	position.Amount = position.Amount.Add(riseAmount)
	if err := k.UserPositions.Set(ctx, collections.Join(msg.Sender, denom), position); err != nil {
		return nil, err
	}

	// Mint rise tokens to user
	riseCoins := sdk.NewCoins(sdk.NewCoin(market.RiseDenom, riseAmount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, riseCoins); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		sdk.AccAddress(sender),
		riseCoins,
	); err != nil {
		return nil, err
	}

	// Emit events
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSupply,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRiseAmount, riseAmount.String()),
		),
	})

	return &types.MsgSupplyResponse{}, nil
}
