package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
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

	// Validate amount
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return nil, errors.Wrap(types.ErrInvalidRequest, "invalid amount")
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

	// Check balance
	denom := GetTokenDenom(msg.TokenCreator)
	balance := k.bankKeeper.GetBalance(ctx, adminAddr, denom)
	if balance.Amount.LT(msg.Amount) {
		return nil, types.ErrInsufficientBalance
	}

	// Create coins to burn
	coins := sdk.NewCoins(sdk.NewCoin(denom, msg.Amount))

	// Send coins from admin to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, adminAddr, types.ModuleName, coins); err != nil {
		return nil, err
	}

	// Burn coins from module
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
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
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
		),
	})

	return &types.MsgBurnResponse{}, nil
}
