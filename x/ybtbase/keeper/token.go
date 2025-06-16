package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

// GetToken returns the token information for a given creator
func (k Keeper) GetToken(ctx context.Context, creator string) (types.Token, bool) {
	token, err := k.Tokens.Get(ctx, creator)
	if err != nil {
		return types.Token{}, false
	}
	return token, true
}

// SetToken sets the token information for a given creator
func (k Keeper) SetToken(ctx context.Context, creator string, token types.Token) error {
	return k.Tokens.Set(ctx, creator, token)
}

// GetGlobalRewardIndex returns the global reward index for a token
func (k Keeper) GetGlobalRewardIndex(ctx context.Context, tokenCreator string) math.LegacyDec {
	index, err := k.GlobalRewardIndex.Get(ctx, tokenCreator)
	if err != nil {
		// Default to 1 if not found
		return math.LegacyOneDec()
	}
	return index
}

// SetGlobalRewardIndex sets the global reward index for a token
func (k Keeper) SetGlobalRewardIndex(ctx context.Context, tokenCreator string, index math.LegacyDec) error {
	return k.GlobalRewardIndex.Set(ctx, tokenCreator, index)
}

// GetUserLastRewardIndex returns the user's last reward index for a token
func (k Keeper) GetUserLastRewardIndex(ctx context.Context, tokenCreator, userAddr string) math.LegacyDec {
	key := collections.Join(tokenCreator, userAddr)
	index, err := k.UserLastRewardIndex.Get(ctx, key)
	if err != nil {
		// Default to current global index if not found
		return k.GetGlobalRewardIndex(ctx, tokenCreator)
	}
	return index
}

// SetUserLastRewardIndex sets the user's last reward index for a token
func (k Keeper) SetUserLastRewardIndex(ctx context.Context, tokenCreator, userAddr string, index math.LegacyDec) error {
	key := collections.Join(tokenCreator, userAddr)
	return k.UserLastRewardIndex.Set(ctx, key, index)
}

// HasPermission checks if a user has permission for a token
func (k Keeper) HasPermission(ctx context.Context, tokenCreator, userAddr string) bool {
	key := collections.Join(tokenCreator, userAddr)
	has, err := k.Permissions.Get(ctx, key)
	if err != nil {
		return false
	}
	return has
}

// SetPermission sets the permission for a user on a token
func (k Keeper) SetPermission(ctx context.Context, tokenCreator, userAddr string, hasPermission bool) error {
	key := collections.Join(tokenCreator, userAddr)
	return k.Permissions.Set(ctx, key, hasPermission)
}

// GetTokenDenom returns the denom for a YBT token
func GetTokenDenom(tokenCreator string) string {
	return types.GetDenom(tokenCreator)
}

// GetYieldPoolAddress returns the module account address for a token's yield pool
func GetYieldPoolAddress(tokenCreator string) sdk.AccAddress {
	return authtypes.NewModuleAddress("ybtbase/yield/" + tokenCreator)
}
