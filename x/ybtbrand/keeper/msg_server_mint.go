package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
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
	if !msg.Ybt.IsValid() || msg.Ybt.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "amount must be positive")
	}

	// Check if YBT denom matches base YBT
	expectedBaseDenom := fmt.Sprintf("ybtbase/%s", token.BaseYbtCreator)
	if msg.Ybt.Denom != expectedBaseDenom {
		return nil, errors.Wrapf(types.ErrInvalidRequest, "invalid base YBT denom: expected %s, got %s", expectedBaseDenom, msg.Ybt.Denom)
	}

	// Transfer base YBT from admin to collateral pool
	collateralAddr := GetCollateralPoolAddress(msg.TokenCreator)
	baseCoins := sdk.NewCoins(msg.Ybt)
	if err := k.bankKeeper.SendCoins(ctx, sdk.AccAddress(adminAddr), collateralAddr, baseCoins); err != nil {
		return nil, err
	}

	// Mint brand tokens to admin
	brandDenom := GetTokenDenom(msg.TokenCreator)
	brandCoins := sdk.NewCoins(sdk.NewCoin(brandDenom, msg.Ybt.Amount))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, brandCoins); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(adminAddr), brandCoins); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Ybt.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, brandDenom),
		),
	})

	return &types.MsgMintResponse{}, nil
}
