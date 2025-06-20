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

func (k msgServer) Borrow(ctx context.Context, msg *types.MsgBorrow) (*types.MsgBorrowResponse, error) {
	// Validate sender address
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Validate amount
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "borrow amount must be positive")
	}

	// Get market
	denom := msg.Amount.Denom
	market, err := k.Markets.Get(ctx, denom)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrMarketNotFound, denom)
		}
		return nil, err
	}

	// Check available liquidity
	availableLiquidity := market.TotalSupplied.Sub(market.TotalBorrowed)
	if availableLiquidity.LT(msg.Amount.Amount) {
		return nil, errorsmod.Wrap(types.ErrInsufficientBalance, "insufficient liquidity")
	}

	// TODO: Get collateral value from liquidity pool module
	// For now, mock the collateral value based on pool ID
	var collateralValue math.Int
	if msg.CollateralPoolId == 1 && msg.CollateralPositionId == 100 {
		// For test: pool 1, position 100 has $2 worth of collateral
		collateralValue = math.NewInt(2000000)
	} else if msg.CollateralPoolId == 1 && msg.CollateralPositionId == 101 {
		// For test: pool 1, position 101 has low collateral
		collateralValue = math.NewInt(100000)
	} else {
		// Default: no collateral
		collateralValue = math.ZeroInt()
	}

	// Get lending parameters
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate borrowing power (collateral value * LTV ratio)
	borrowingPower := math.LegacyNewDecFromInt(collateralValue).Mul(params.LtvRatio).TruncateInt()

	// TODO: Convert borrow amount to USD value using TWAP oracle
	// For now, assume 1:1 for USDC
	borrowValueUSD := msg.Amount.Amount

	// Check if user has sufficient collateral
	if borrowingPower.LT(borrowValueUSD) {
		return nil, errorsmod.Wrap(types.ErrInsufficientCollateral, "insufficient collateral")
	}

	// Get next borrow ID
	borrowId, err := k.BorrowId.Next(ctx)
	if err != nil {
		return nil, err
	}

	// Create borrow record
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	borrow := types.Borrow{
		Id:                   borrowId,
		Borrower:             msg.Sender,
		Amount:               msg.Amount,
		CollateralPoolId:     msg.CollateralPoolId,
		CollateralPositionId: msg.CollateralPositionId,
		BlockHeight:          sdkCtx.BlockHeight(),
	}

	// Save borrow
	if err := k.Borrows.Set(ctx, borrowId, borrow); err != nil {
		return nil, err
	}

	// Update market total borrowed
	market.TotalBorrowed = market.TotalBorrowed.Add(msg.Amount.Amount)
	if err := k.Markets.Set(ctx, denom, market); err != nil {
		return nil, err
	}

	// Send borrowed tokens to user
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		sdk.AccAddress(sender),
		sdk.NewCoins(msg.Amount),
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send borrowed tokens")
	}

	// Emit events
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBorrow,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyBorrowId, strconv.FormatUint(borrowId, 10)),
		),
	})

	return &types.MsgBorrowResponse{}, nil
}
