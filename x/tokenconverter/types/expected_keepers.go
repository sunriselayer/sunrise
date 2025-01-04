package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	banktypes "cosmossdk.io/x/bank/types"
)

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here

	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error

	GetSupply(ctx context.Context, denom string) sdk.Coin
	HasSupply(ctx context.Context, denom string) bool
	// GetPaginatedTotalSupply(ctx context.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error)
	IterateTotalSupply(ctx context.Context, cb func(sdk.Coin) bool)
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	HasDenomMetaData(ctx context.Context, denom string) bool
	SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata)
	GetAllDenomMetaData(ctx context.Context) []banktypes.Metadata
	IterateAllDenomMetaData(ctx context.Context, cb func(banktypes.Metadata) bool)

	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error

	DelegateCoins(ctx context.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoins(ctx context.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error
}
