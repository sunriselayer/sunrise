package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) Send(ctx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	// Validate addresses
	fromAddr, err := k.addressCodec.StringToBytes(msg.FromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "invalid from address")
	}

	toAddr, err := k.addressCodec.StringToBytes(msg.ToAddress)
	if err != nil {
		return nil, errors.Wrap(err, "invalid to address")
	}

	// Get the token
	token, found := k.Keeper.GetToken(ctx, msg.TokenCreator)
	if !found {
		return nil, types.ErrTokenNotFound
	}

	// Check permission based on permission mode
	switch token.PermissionMode {
	case types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS:
		// No restrictions, anyone can send
	case types.PermissionMode_PERMISSION_MODE_WHITELIST:
		// Check if both sender and receiver are whitelisted
		fromAllowed, _ := k.Keeper.Permissions.Get(ctx, collections.Join(msg.TokenCreator, msg.FromAddress))
		toAllowed, _ := k.Keeper.Permissions.Get(ctx, collections.Join(msg.TokenCreator, msg.ToAddress))
		if !fromAllowed || !toAllowed {
			return nil, errors.Wrap(types.ErrUnauthorized, "address not whitelisted")
		}
	case types.PermissionMode_PERMISSION_MODE_BLACKLIST:
		// Check if either sender or receiver is blacklisted
		fromBlacklisted, _ := k.Keeper.Permissions.Get(ctx, collections.Join(msg.TokenCreator, msg.FromAddress))
		toBlacklisted, _ := k.Keeper.Permissions.Get(ctx, collections.Join(msg.TokenCreator, msg.ToAddress))
		if fromBlacklisted || toBlacklisted {
			return nil, errors.Wrap(types.ErrUnauthorized, "address blacklisted")
		}
	default:
		return nil, errors.Wrap(types.ErrInvalidPermissionMode, "invalid permission mode")
	}

	// Get the denom
	denom := types.GetDenom(msg.TokenCreator)

	// Create coins
	amount := msg.Amount
	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Perform the transfer using bank keeper
	if err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coins); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSend,
			sdk.NewAttribute(types.AttributeKeyFrom, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyTo, msg.ToAddress),
			sdk.NewAttribute(types.AttributeKeyTokenCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return &types.MsgSendResponse{}, nil
}