package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) ClaimYield(ctx context.Context, msg *types.MsgClaimYield) (*types.MsgClaimYieldResponse, error) {
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

	// Check permission based on permission mode
	switch token.PermissionMode {
	case types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS:
		// No restrictions
	case types.PermissionMode_PERMISSION_MODE_WHITELIST:
		// Check if sender is whitelisted
		allowed, err := k.Keeper.Permissions.Get(ctx, collections.Join(msg.TokenCreator, msg.Sender))
		if err != nil || !allowed {
			return nil, types.ErrNoPermission
		}
	case types.PermissionMode_PERMISSION_MODE_BLACKLIST:
		// Check if sender is blacklisted
		blacklisted, _ := k.Keeper.Permissions.Get(ctx, collections.Join(msg.TokenCreator, msg.Sender))
		if blacklisted {
			return nil, types.ErrNoPermission
		}
	default:
		return nil, errors.Wrap(types.ErrInvalidRequest, "invalid permission mode")
	}

	// Get user balance
	denom := GetTokenDenom(msg.TokenCreator)
	balance := k.bankKeeper.GetBalance(ctx, senderAddr, denom)
	if balance.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no balance")
	}

	// Get global reward index and user's last reward index
	globalIndex := k.Keeper.GetGlobalRewardIndex(ctx, msg.TokenCreator)
	userLastIndex := k.Keeper.GetUserLastRewardIndex(ctx, msg.TokenCreator, msg.Sender)

	// Calculate yield amount
	// yield = balance * (globalIndex - userLastIndex)
	indexDiff := globalIndex.Sub(userLastIndex)
	if indexDiff.IsZero() || indexDiff.IsNegative() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no yield to claim")
	}

	yieldAmount := indexDiff.MulInt(balance.Amount).TruncateInt()
	if yieldAmount.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no yield to claim")
	}

	// Check yield pool balance
	yieldPoolAddr := GetYieldPoolAddress(msg.TokenCreator)
	yieldPoolBalance := k.bankKeeper.GetBalance(ctx, yieldPoolAddr, denom)
	if yieldPoolBalance.Amount.LT(yieldAmount) {
		return nil, errors.Wrap(types.ErrInsufficientBalance, "insufficient yield pool balance")
	}

	// Transfer yield from pool to user
	yieldCoins := sdk.NewCoins(sdk.NewCoin(denom, yieldAmount))
	if err := k.bankKeeper.SendCoins(ctx, yieldPoolAddr, senderAddr, yieldCoins); err != nil {
		return nil, err
	}

	// Update user's last reward index
	if err := k.Keeper.SetUserLastRewardIndex(ctx, msg.TokenCreator, msg.Sender, globalIndex); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimYield,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyClaimer, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyYieldAmount, yieldAmount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
		),
	})

	return &types.MsgClaimYieldResponse{}, nil
}
