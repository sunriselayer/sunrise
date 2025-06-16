package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) GrantPermission(ctx context.Context, msg *types.MsgGrantPermission) (*types.MsgGrantPermissionResponse, error) {
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

	// For whitelist mode, grant permission means allowing access
	// For blacklist mode, we don't grant permissions (only revoke to unblock)
	if token.PermissionMode != types.PermissionMode_PERMISSION_MODE_WHITELIST {
		return nil, errors.Wrap(types.ErrInvalidRequest, "can only grant permissions in whitelist mode")
	}

	// Grant permission (set to true for whitelist)
	if err := k.Keeper.Permissions.Set(ctx, collections.Join(msg.TokenCreator, msg.Target), true); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeGrantPermission,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyTarget, msg.Target),
		),
	})

	return &types.MsgGrantPermissionResponse{}, nil
}