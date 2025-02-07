package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here

	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	IsSendEnabledCoins(ctx context.Context, coins ...sdk.Coin) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// TransferKeeper defines the expected interface for the IBC Transfer module.
type TransferKeeper interface {
	Transfer(ctx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error)
	GetTotalEscrowForDenom(ctx context.Context, denom string) sdk.Coin
	SetTotalEscrowForDenom(ctx context.Context, coin sdk.Coin)
}

// LiquidityPoolKeeper defines the expected interface for the liquidity pool module.
type LiquidityPoolKeeper interface {
	GetPool(ctx context.Context, id uint64) (val lptypes.Pool, found bool)
	CalculateResultExactAmountIn(ctx sdk.Context, pool lptypes.Pool, tokenIn sdk.Coin, denomOut string, feeEnabled bool) (amountOut math.Int, err error)
	CalculateResultExactAmountOut(ctx sdk.Context, pool lptypes.Pool, tokenOut sdk.Coin, denomIn string, feeEnabled bool) (amountIn math.Int, err error)
	SwapExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, pool lptypes.Pool, tokenIn sdk.Coin, denomOut string, feeEnabled bool) (amountOut math.Int, err error)
	SwapExactAmountOut(ctx sdk.Context, sender sdk.AccAddress, pool lptypes.Pool, tokenOut sdk.Coin, denomIn string, feeEnabled bool) (amountIn math.Int, err error)
}
