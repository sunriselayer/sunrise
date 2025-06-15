package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) ClaimYields(ctx context.Context, msg *types.MsgClaimYields) (*types.MsgClaimYieldsResponse, error) {
	// Validate sender address
	senderAddr, err := k.addressCodec.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errors.Wrap(err, "invalid sender address")
	}

	// Validate token creator address
	if _, err := k.addressCodec.StringToBytes(msg.TokenCreator); err != nil {
		return nil, errors.Wrap(err, "invalid token creator address")
	}

	// Get token
	token, found := k.Keeper.GetToken(ctx, msg.TokenCreator)
	if !found {
		return nil, types.ErrTokenNotFound
	}

	// Validate denoms
	if len(msg.Denoms) == 0 {
		return nil, errors.Wrap(types.ErrInvalidRequest, "denoms cannot be empty")
	}

	// Get user's brand token balance
	brandDenom := GetTokenDenom(msg.TokenCreator)
	balance := k.bankKeeper.GetBalance(ctx, sdk.AccAddress(senderAddr), brandDenom)
	if balance.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no balance")
	}

	// Process each yield denom
	totalClaimed := sdk.NewCoins()
	for _, denom := range msg.Denoms {
		// Get global yield index
		globalIndex, found := k.Keeper.GetYieldIndex(ctx, msg.TokenCreator, denom)
		if !found {
			continue // Skip if no yield index exists
		}

		// Get user's last yield index
		userLastIndex, _ := k.Keeper.GetUserLastYieldIndex(ctx, msg.TokenCreator, msg.Sender, denom)
		// If not found, defaults to 1.0

		// Calculate yield amount
		// yield = balance * (globalIndex - userLastIndex)
		indexDiff := globalIndex.Sub(userLastIndex)
		if indexDiff.IsZero() || indexDiff.IsNegative() {
			continue // No yield to claim for this denom
		}

		yieldAmount := indexDiff.MulInt(balance.Amount).TruncateInt()
		if yieldAmount.IsZero() {
			continue
		}

		// Check yield pool balance
		yieldPoolAddr := GetYieldPoolAddress(msg.TokenCreator, denom)
		yieldPoolBalance := k.bankKeeper.GetBalance(ctx, yieldPoolAddr, denom)
		if yieldPoolBalance.Amount.LT(yieldAmount) {
			return nil, errors.Wrapf(types.ErrInsufficientBalance, "insufficient yield pool balance for %s", denom)
		}

		// Transfer yield from pool to user
		yieldCoins := sdk.NewCoins(sdk.NewCoin(denom, yieldAmount))
		if err := k.bankKeeper.SendCoins(ctx, yieldPoolAddr, sdk.AccAddress(senderAddr), yieldCoins); err != nil {
			return nil, err
		}

		// Update user's last yield index
		if err := k.Keeper.SetUserLastYieldIndex(ctx, msg.TokenCreator, msg.Sender, denom, globalIndex); err != nil {
			return nil, err
		}

		totalClaimed = totalClaimed.Add(yieldCoins...)
	}

	if totalClaimed.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no yield to claim")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimYields,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyClaimer, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyYieldAmount, totalClaimed.String()),
		),
	})

	// Log for debugging
	k.Logger(ctx).Info("Claimed yields",
		"claimer", msg.Sender,
		"token_creator", msg.TokenCreator,
		"claimed", totalClaimed.String(),
		"token", token.BaseYbtCreator,
	)

	return &types.MsgClaimYieldsResponse{}, nil
}
