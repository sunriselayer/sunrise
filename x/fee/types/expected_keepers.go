package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// LiquidityPoolKeeper defines the expected interface for the liquidity pool module.
type LiquidityPoolKeeper interface {
	GetPool(ctx sdk.Context, id uint64) (val lptypes.Pool, found bool, err error)
	SwapExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, pool lptypes.Pool, tokenIn sdk.Coin, denomOut string, feeEnabled bool) (amountOut math.Int, err error)
}
