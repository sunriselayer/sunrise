package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) GrantYieldPermission(ctx context.Context, msg *types.MsgGrantYieldPermission) (*types.MsgGrantYieldPermissionResponse, error) {
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

	// Check if token is permissioned
	if !token.Permissioned {
		return nil, errors.Wrap(types.ErrInvalidRequest, "token is not permissioned")
	}

	// Check if msg sender is admin
	if token.Admin != msg.Admin {
		return nil, types.ErrUnauthorized
	}

	// Grant permission
	if err := k.Keeper.SetYieldPermission(ctx, msg.TokenCreator, msg.Target, true); err != nil {
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

	return &types.MsgGrantYieldPermissionResponse{}, nil
}
