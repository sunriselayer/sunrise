package types

//go:generate mockgen -destination ../testutil/expected_keepers_mocks.go -package testutil . AuthKeeper,BankKeeper,YbtbaseKeeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ybtbasetypes "github.com/sunriselayer/sunrise/x/ybtbase/types"
)

// AuthKeeper defines the expected auth keeper
type AuthKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI)
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}

// YbtbaseKeeper defines the expected ybtbase keeper
type YbtbaseKeeper interface {
	GetToken(ctx context.Context, creator string) (ybtbasetypes.Token, bool)
	GetGlobalRewardIndex(ctx context.Context, creator string) math.LegacyDec
	GetUserLastRewardIndex(ctx context.Context, creator, user string) math.LegacyDec
	SetUserLastRewardIndex(ctx context.Context, creator, user string, index math.LegacyDec) error
	HasPermission(ctx context.Context, creator, user string) bool
}
