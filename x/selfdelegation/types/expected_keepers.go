package types

import (
	"context"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountsKeeper interface {
	Init(
		ctx context.Context,
		accountType string,
		creator []byte,
		initRequest transaction.Msg,
		funds sdk.Coins,
		addressSeed []byte,
	) (transaction.Msg, []byte, error)

	Query(
		ctx context.Context,
		accountAddr []byte,
		queryRequest transaction.Msg,
	) (transaction.Msg, error)
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here

	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
}

type FeeKeeper interface {
	FeeDenom(ctx context.Context) (string, error)
}

type TokenConverterKeeper interface {
	Convert(ctx context.Context, amount math.Int, address sdk.AccAddress) error
	ConvertReverse(ctx context.Context, amount math.Int, address sdk.AccAddress) error
}
