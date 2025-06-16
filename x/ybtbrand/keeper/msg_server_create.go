package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
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

	// Validate base YBT creator address
	if _, err := k.addressCodec.StringToBytes(msg.BaseYbtCreator); err != nil {
		return nil, errors.Wrap(err, "invalid base YBT creator address")
	}

	// Check if token already exists
	if k.Keeper.HasToken(ctx, msg.Creator) {
		return nil, errors.Wrap(types.ErrTokenAlreadyExists, "token already exists")
	}

	// Check if base YBT token exists
	baseToken, found := k.ybtbaseKeeper.GetToken(ctx, msg.BaseYbtCreator)
	if !found {
		return nil, errors.Wrap(types.ErrTokenNotFound, "base YBT token not found")
	}

	// Create token
	token := types.Token{
		Creator:        msg.Creator,
		Admin:          msg.Admin,
		BaseYbtCreator: msg.BaseYbtCreator,
	}

	if err := k.Keeper.SetToken(ctx, msg.Creator, token); err != nil {
		return nil, err
	}

	// Create collateral pool module account
	collateralAddr := GetCollateralPoolAddress(msg.Creator)
	if acc := k.authKeeper.GetAccount(ctx, collateralAddr); acc == nil {
		acc := authtypes.NewModuleAccount(
			authtypes.NewBaseAccountWithAddress(collateralAddr),
			collateralAddr.String(),
		)
		k.authKeeper.SetModuleAccount(ctx, acc)
	}

	// Check if initial supply is zero
	denom := GetTokenDenom(msg.Creator)
	supply := k.bankKeeper.GetSupply(ctx, denom)
	if !supply.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "token already has supply")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreate,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyBaseYbtCreator, msg.BaseYbtCreator),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
		),
	})

	// Log base token info for debugging
	k.Logger(ctx).Info("Created ybtbrand token",
		"creator", msg.Creator,
		"admin", msg.Admin,
		"base_ybt_creator", msg.BaseYbtCreator,
		"base_token_admin", baseToken.Admin,
	)

	return &types.MsgCreateResponse{}, nil
}
