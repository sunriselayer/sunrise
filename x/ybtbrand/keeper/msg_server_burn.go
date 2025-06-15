package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	// Validate admin address
	adminAddr, err := k.addressCodec.StringToBytes(msg.Admin)
	if err != nil {
		return nil, errors.Wrap(err, "invalid admin address")
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

	// Check if sender is admin
	if token.Admin != msg.Admin {
		return nil, types.ErrUnauthorized
	}

	// Validate amount
	if !msg.Amount.IsPositive() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "amount must be positive")
	}

	// Burn brand tokens from admin
	brandDenom := GetTokenDenom(msg.TokenCreator)
	brandCoins := sdk.NewCoins(sdk.NewCoin(brandDenom, msg.Amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(adminAddr), types.ModuleName, brandCoins); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, brandCoins); err != nil {
		return nil, err
	}

	// Transfer base YBT from collateral pool to admin
	collateralAddr := GetCollateralPoolAddress(msg.TokenCreator)
	baseDenom := fmt.Sprintf("ybtbase/%s", token.BaseYbtCreator)
	baseCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, msg.Amount))
	if err := k.bankKeeper.SendCoins(ctx, collateralAddr, sdk.AccAddress(adminAddr), baseCoins); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, brandDenom),
		),
	})

	return &types.MsgBurnResponse{}, nil
}