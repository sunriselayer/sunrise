package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
	swaptypes "github.com/sunriselayer/sunrise/x/swap/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// LiquidityPoolKeeper defines the expected interface for the liquidity pool module.
type LiquidityPoolKeeper interface {
	GetPool(ctx context.Context, id uint64) (val lptypes.Pool, found bool, err error)
}

// SwapKeeper defines the expected interface for the swap module.
type SwapKeeper interface {
	SwapExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, interfaceProvider string, route swaptypes.Route, amountIn math.Int, minAmountOut math.Int) (result swaptypes.RouteResult, interfaceFee math.Int, err error)
}
