package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) RevokePermission(ctx context.Context, msg *types.MsgRevokePermission) (*types.MsgRevokePermissionResponse, error) {
	// Validate admin address
	if _, err := k.addressCodec.StringToBytes(msg.Admin); err != nil {
		return nil, errors.Wrap(err, "invalid admin address")
	}

	// Validate token creator address
	if _, err := k.addressCodec.StringToBytes(msg.TokenCreator); err != nil {
		return nil, errors.Wrap(err, "invalid token creator address")
	}

	// Validate target address
	if _, err := k.addressCodec.StringToBytes(msg.Target); err != nil {
		return nil, errors.Wrap(err, "invalid target address")
	}

	// Get token
	token, found := k.Keeper.GetToken(ctx, msg.TokenCreator)
	if !found {
		return nil, types.ErrTokenNotFound
	}

	// Check if msg sender is admin
	if token.Admin != msg.Admin {
		return nil, types.ErrUnauthorized
	}

	// Handle based on permission mode
	switch token.PermissionMode {
	case types.PermissionMode_PERMISSION_MODE_WHITELIST:
		// For whitelist, revoke means removing permission (delete from map)
		if err := k.Keeper.Permissions.Remove(ctx, collections.Join(msg.TokenCreator, msg.Target)); err != nil {
			return nil, err
		}
		
		// Check if there are any remaining permissions
		hasAnyPermissions := false
		iter, err := k.Keeper.Permissions.Iterate(ctx, collections.NewPrefixedPairRange[string, string](msg.TokenCreator))
		if err != nil {
			return nil, err
		}
		defer iter.Close()
		
		if iter.Valid() {
			hasAnyPermissions = true
		}
		
		// Update send enabled status based on whether any permissions remain
		denom := types.GetDenom(msg.TokenCreator)
		k.bankKeeper.SetSendEnabled(ctx, denom, hasAnyPermissions)
		
	case types.PermissionMode_PERMISSION_MODE_BLACKLIST:
		// For blacklist, revoke means adding to blacklist (set to true)
		if err := k.Keeper.Permissions.Set(ctx, collections.Join(msg.TokenCreator, msg.Target), true); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Wrap(types.ErrInvalidRequest, "cannot revoke permissions in permissionless mode")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevokePermission,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyTarget, msg.Target),
		),
	})

	return &types.MsgRevokePermissionResponse{}, nil
}