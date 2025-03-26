package types

import (
	"context"

	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "cosmossdk.io/x/staking/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here

	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin

	SetSendEnabled(ctx context.Context, denom string, value bool)

	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	IsSendEnabledCoins(ctx context.Context, coins ...sdk.Coin) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, address []byte, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type StakingKeeper interface {
	ValidatorAddressCodec() addresscodec.Codec

	GetParams(ctx context.Context) (stakingtypes.Params, error)
	BondDenom(ctx context.Context) (string, error)
	IterateDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress, cb func(delegation stakingtypes.Delegation) (stop bool)) error
}

type FeeKeeper interface {
	FeeDenom(ctx context.Context) (string, error)
}

type TokenConverterKeeper interface {
	Convert(ctx context.Context, amount math.Int, address sdk.AccAddress) error
	ConvertReverse(ctx context.Context, amount math.Int, address sdk.AccAddress) error
}
