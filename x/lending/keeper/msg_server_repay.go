package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	"github.com/sunriselayer/sunrise/x/lending/types"
)

func (k msgServer) Repay(ctx context.Context, msg *types.MsgRepay) (*types.MsgRepayResponse, error) {
	// Validate sender address
	sender, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Validate amount
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "repayment amount must be positive")
	}

	// Get borrow
	borrow, err := k.Borrows.Get(ctx, msg.BorrowId)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrBorrowNotFound, "borrow not found")
		}
		return nil, err
	}

	// Check if sender is the borrower
	if borrow.Borrower != msg.Sender {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only borrower can repay")
	}

	// Check denom matches
	if borrow.Amount.Denom != msg.Amount.Denom {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "repayment denom mismatch")
	}

	// Check repayment amount doesn't exceed debt
	if msg.Amount.Amount.GT(borrow.Amount.Amount) {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "repayment exceeds debt")
	}

	// Transfer repayment from user to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.AccAddress(sender),
		types.ModuleName,
		sdk.NewCoins(msg.Amount),
	); err != nil {
		return nil, errorsmod.Wrap(err, "failed to transfer repayment")
	}

	// Update market total borrowed
	market, err := k.Markets.Get(ctx, borrow.Amount.Denom)
	if err != nil {
		return nil, err
	}
	market.TotalBorrowed = market.TotalBorrowed.Sub(msg.Amount.Amount)
	if err := k.Markets.Set(ctx, borrow.Amount.Denom, market); err != nil {
		return nil, err
	}

	// Update or remove borrow
	if msg.Amount.Amount.Equal(borrow.Amount.Amount) {
		// Full repayment - remove borrow
		if err := k.Borrows.Remove(ctx, msg.BorrowId); err != nil {
			return nil, err
		}
		
		// TODO: Release collateral from liquidity pool
	} else {
		// Partial repayment - update borrow amount
		borrow.Amount.Amount = borrow.Amount.Amount.Sub(msg.Amount.Amount)
		if err := k.Borrows.Set(ctx, msg.BorrowId, borrow); err != nil {
			return nil, err
		}
	}

	// Emit events
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRepay,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyBorrowId, strconv.FormatUint(msg.BorrowId, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
	})

	return &types.MsgRepayResponse{}, nil
}
