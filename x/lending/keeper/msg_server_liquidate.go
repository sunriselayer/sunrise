package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func (k msgServer) Liquidate(ctx context.Context, msg *types.MsgLiquidate) (*types.MsgLiquidateResponse, error) {
	// Validate sender address
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Validate amount
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "liquidation amount must be positive")
	}

	// Get borrow
	borrow, err := k.Borrows.Get(ctx, msg.BorrowId)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrBorrowNotFound, "borrow not found")
		}
		return nil, err
	}

	// Check denom matches
	if borrow.Amount.Denom != msg.Amount.Denom {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "liquidation denom mismatch")
	}

	// Check liquidation amount doesn't exceed debt
	if msg.Amount.Amount.GT(borrow.Amount.Amount) {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "liquidation amount exceeds debt")
	}

	// Get market
	market, err := k.Markets.Get(ctx, borrow.Amount.Denom)
	if err != nil {
		return nil, err
	}

	// Get lending parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: Get collateral value from liquidity pool module
	// For now, mock the collateral value based on position
	var collateralValue math.Int
	if borrow.CollateralPoolId == 1 && borrow.CollateralPositionId == 100 {
		// Healthy position
		collateralValue = math.NewInt(2000000) // $2 worth
	} else if borrow.CollateralPoolId == 1 && borrow.CollateralPositionId == 200 {
		// Undercollateralized position
		collateralValue = math.NewInt(800000) // $0.8 worth
	} else {
		collateralValue = math.ZeroInt()
	}

	// TODO: Convert borrow amount to USD value using TWAP oracle
	// For now, assume 1:1 for USDC
	borrowValueUSD := borrow.Amount.Amount

	// Calculate health factor (collateral value / borrow value)
	healthFactor := math.LegacyNewDecFromInt(collateralValue).Quo(math.LegacyNewDecFromInt(borrowValueUSD))

	// Check if position is undercollateralized
	if healthFactor.GTE(params.LiquidationThreshold) {
		return nil, errorsmod.Wrap(types.ErrUndercollateralized, "position is not undercollateralized")
	}

	// Transfer liquidation payment from liquidator to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.AccAddress(sender),
		types.ModuleName,
		sdk.NewCoins(msg.Amount),
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to transfer liquidation payment")
	}

	// Update market total borrowed
	market.TotalBorrowed = market.TotalBorrowed.Sub(msg.Amount.Amount)
	if err := k.Markets.Set(ctx, borrow.Amount.Denom, market); err != nil {
		return nil, err
	}

	// Update or remove borrow
	if msg.Amount.Amount.Equal(borrow.Amount.Amount) {
		// Full liquidation - remove borrow
		if err := k.Borrows.Remove(ctx, msg.BorrowId); err != nil {
			return nil, err
		}
	} else {
		// Partial liquidation - update borrow amount
		borrow.Amount.Amount = borrow.Amount.Amount.Sub(msg.Amount.Amount)
		if err := k.Borrows.Set(ctx, msg.BorrowId, borrow); err != nil {
			return nil, err
		}
	}

	// TODO: Calculate and transfer liquidation reward to liquidator
	// This would involve:
	// 1. Calculating liquidation bonus (e.g., 5% of liquidated amount)
	// 2. Claiming proportional collateral from liquidity pool
	// 3. Transferring collateral to liquidator

	// Emit events
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLiquidate,
			sdk.NewAttribute(types.AttributeKeyLiquidator, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyBorrower, borrow.Borrower),
			sdk.NewAttribute(types.AttributeKeyBorrowId, strconv.FormatUint(msg.BorrowId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
	})

	return &types.MsgLiquidateResponse{}, nil
}
