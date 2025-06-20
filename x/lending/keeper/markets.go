package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

// GetMarket retrieves a market by denom
func (k Keeper) GetMarket(ctx context.Context, denom string) (types.Market, error) {
	return k.Markets.Get(ctx, denom)
}

// SetMarket stores a market
func (k Keeper) SetMarket(ctx context.Context, market types.Market) error {
	return k.Markets.Set(ctx, market.Denom, market)
}

// GetUserPosition retrieves a user's position in a market
func (k Keeper) GetUserPosition(ctx context.Context, userAddress, denom string) (types.UserPosition, error) {
	return k.UserPositions.Get(ctx, collections.Join(userAddress, denom))
}

// SetUserPosition stores a user's position
func (k Keeper) SetUserPosition(ctx context.Context, position types.UserPosition) error {
	return k.UserPositions.Set(ctx, collections.Join(position.UserAddress, position.Denom), position)
}