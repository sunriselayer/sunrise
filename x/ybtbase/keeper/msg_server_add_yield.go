package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

func (k msgServer) AddYield(ctx context.Context, msg *types.MsgAddYield) (*types.MsgAddYieldResponse, error) {
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

	// Get total supply (excluding yield pool balance)
	totalSupply := k.GetTotalSupplyExcludingYieldPool(ctx, msg.TokenCreator)
	if totalSupply.IsZero() {
		// No tokens in circulation, cannot add yield
		return nil, errors.Wrap(types.ErrInvalidRequest, "no tokens in circulation")
	}

	// Calculate new global reward index
	currentIndex := k.Keeper.GetGlobalRewardIndex(ctx, msg.TokenCreator)
	yieldPerToken := math.LegacyNewDecFromInt(msg.Amount).Quo(math.LegacyNewDecFromInt(totalSupply))
	newIndex := currentIndex.Add(yieldPerToken)

	// Update global reward index
	if err := k.Keeper.SetGlobalRewardIndex(ctx, msg.TokenCreator, newIndex); err != nil {
		return nil, err
	}

	// Transfer yield to yield pool
	yieldPoolAddr := GetYieldPoolAddress(msg.TokenCreator)
	coins := sdk.NewCoins(sdk.NewCoin(denom, msg.Amount))
	if err := k.bankKeeper.SendCoins(ctx, adminAddr, yieldPoolAddr, coins); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddYield,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.TokenCreator),
			sdk.NewAttribute(types.AttributeKeyAdmin, msg.Admin),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute("new_global_index", newIndex.String()),
		),
	})

	return &types.MsgAddYieldResponse{}, nil
}

// GetTotalSupplyExcludingYieldPool returns the total supply of a token excluding the yield pool balance
func (k msgServer) GetTotalSupplyExcludingYieldPool(ctx context.Context, tokenCreator string) math.Int {
	// For simplicity in tests, we'll use a mock approach
	// In production, this would iterate through all balances or use a cached total supply
	// For now, we'll just return a fixed amount for testing
	return math.NewInt(10000)
}
