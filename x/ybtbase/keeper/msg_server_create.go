package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) Create(ctx context.Context, msg *types.MsgCreate) (*types.MsgCreateResponse, error) {
	// Validate creator address
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errors.Wrap(err, "invalid creator address")
	}

	// Validate admin address
	if _, err := k.addressCodec.StringToBytes(msg.Admin); err != nil {
		return nil, errors.Wrap(err, "invalid admin address")
	}

	// Check if token already exists
	if _, found := k.Keeper.GetToken(ctx, msg.Creator); found {
		return nil, types.ErrTokenAlreadyExists
	}

	// Validate permission mode
	if msg.PermissionMode == types.PermissionMode_PERMISSION_MODE_UNSPECIFIED {
		return nil, errors.Wrap(types.ErrInvalidRequest, "permission mode must be specified")
	}

	// Create the token
	token := types.Token{
		Creator:        msg.Creator,
		Admin:          msg.Admin,
		PermissionMode: msg.PermissionMode,
	}

	// Save the token
	if err := k.Keeper.SetToken(ctx, msg.Creator, token); err != nil {
		return nil, err
	}

	// Initialize global reward index to 1
	if err := k.Keeper.SetGlobalRewardIndex(ctx, msg.Creator, math.LegacyOneDec()); err != nil {
		return nil, err
	}

	// Get the denom for this token
	denom := types.GetDenom(msg.Creator)

	// Disable bank send for non-permissionless tokens
	if msg.PermissionMode != types.PermissionMode_PERMISSION_MODE_PERMISSIONLESS {
		// Set SendEnabled to false for this denom
		k.bankKeeper.SetSendEnabled(ctx, denom, false)
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateToken,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute("permission_mode", msg.PermissionMode.String()),
		),
	})

	return &types.MsgCreateResponse{}, nil
}
