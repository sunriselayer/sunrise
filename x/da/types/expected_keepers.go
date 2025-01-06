package types

import (
	"context"

	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// StakingKeeper is expected keeper for slashing module
type SlashingKeeper interface {
	Slash(context.Context, sdk.ConsAddress, math.LegacyDec, int64, int64) error
	Jail(context.Context, sdk.ConsAddress) error
}

// StakingKeeper is expected keeper for staking module
type StakingKeeper interface {
	Validator(ctx context.Context, address sdk.ValAddress) (sdk.ValidatorI, error)
	PowerReduction(ctx context.Context) (res math.Int)
	ValidatorsPowerStoreIterator(ctx context.Context) (corestore.Iterator, error)
	TotalBondedTokens(ctx context.Context) (math.Int, error)
}
