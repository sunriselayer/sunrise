package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

func (k msgServer) ClaimCollateralYield(ctx context.Context, msg *types.MsgClaimCollateralYield) (*types.MsgClaimCollateralYieldResponse, error) {
	// Validate admin address
	adminAddr, err := k.addressCodec.StringToBytes(msg.Admin)
	if err != nil {
		return nil, errors.Wrap(err, "invalid admin address")
	}

	// Validate token creator address
	if _, err := k.addressCodec.StringToBytes(msg.TokenCreator); err != nil {
		return nil, errors.Wrap(err, "invalid token creator address")
	}

	// Validate base YBT creator address
	if _, err := k.addressCodec.StringToBytes(msg.BaseYbtCreator); err != nil {
		return nil, errors.Wrap(err, "invalid base YBT creator address")
	}

	// Get brand token
	token, found := k.Keeper.GetToken(ctx, msg.TokenCreator)
	if !found {
		return nil, types.ErrTokenNotFound
	}

	// Check if sender is admin
	if token.Admin != msg.Admin {
		return nil, types.ErrUnauthorized
	}

	// Check if base YBT creator matches
	if token.BaseYbtCreator != msg.BaseYbtCreator {
		return nil, errors.Wrap(types.ErrInvalidRequest, "base YBT creator mismatch")
	}

	// Get base YBT token info
	baseToken, found := k.ybtbaseKeeper.GetToken(ctx, msg.BaseYbtCreator)
	if !found {
		return nil, errors.Wrap(types.ErrInvalidRequest, "base YBT token not found")
	}

	// For permissioned base YBT, check if admin has yield permission
	if baseToken.Permissioned {
		if !k.ybtbaseKeeper.HasYieldPermission(ctx, msg.BaseYbtCreator, msg.Admin) {
			return nil, errors.Wrap(types.ErrUnauthorized, "no yield permission")
		}
	}

	// Get collateral pool address
	collateralAddr := GetCollateralPoolAddress(msg.TokenCreator)
	baseDenom := types.GetBaseYbtTokenDenom(msg.BaseYbtCreator)

	// Get collateral balance
	collateralBalance := k.bankKeeper.GetBalance(ctx, collateralAddr, baseDenom)
	if collateralBalance.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no collateral balance")
	}

	// Get global reward index
	globalIndex := k.ybtbaseKeeper.GetGlobalRewardIndex(ctx, msg.BaseYbtCreator)

	// Get collateral pool's last reward index
	lastIndex := k.ybtbaseKeeper.GetUserLastRewardIndex(ctx, msg.BaseYbtCreator, collateralAddr.String())

	// Calculate yield amount
	// yield = balance * (globalIndex - lastIndex)
	indexDiff := globalIndex.Sub(lastIndex)
	if indexDiff.IsZero() || indexDiff.IsNegative() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no yield to claim")
	}

	yieldAmount := indexDiff.MulInt(collateralBalance.Amount).TruncateInt()
	if yieldAmount.IsZero() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no yield to claim")
	}

	// Transfer yield from base YBT yield pool to admin
	yieldPoolAddr := types.GetBaseYbtYieldPoolAddress(msg.BaseYbtCreator)
	yieldCoins := sdk.NewCoins(sdk.NewCoin(baseDenom, yieldAmount))
	if err := k.bankKeeper.SendCoins(ctx, yieldPoolAddr, sdk.AccAddress(adminAddr), yieldCoins); err != nil {
		return nil, err
	}

	// Update collateral pool's last reward index
	if err := k.ybtbaseKeeper.SetUserLastRewardIndex(ctx, msg.BaseYbtCreator, collateralAddr.String(), globalIndex); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimCollateralYield,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyBaseYbtCreator, msg.BaseYbtCreator),
			sdk.NewAttribute(types.AttributeKeyYieldAmount, yieldCoins.String()),
		),
	})

	return &types.MsgClaimCollateralYieldResponse{}, nil
}