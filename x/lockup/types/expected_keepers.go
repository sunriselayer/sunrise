package types

import (
	"context"
	"time"

	addresscodec "cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	AddressCodec() addresscodec.Codec
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	ValidatorAddressCodec() addresscodec.Codec
}

type FeeKeeper interface {
	FeeDenom(ctx context.Context) (string, error)
}

type ShareclassKeeper interface {
	Delegate(ctx context.Context, sender sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) (share, rewards sdk.Coins, err error)
	Undelegate(ctx context.Context, sender sdk.AccAddress, recipient sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) (output sdk.Coin, rewards sdk.Coins, CompletionTime time.Time, err error)
	ClaimRewards(ctx context.Context, sender sdk.AccAddress, validatorAddr sdk.ValAddress) (sdk.Coins, error)
}
