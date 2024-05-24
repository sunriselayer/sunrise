package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lptypes "github.com/sunriselayer/sunrise/x/liquiditypool/types"
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
}

// SwapKeeper defines the expected interface for the liquidity module.
type SwapKeeper interface {
	GetPool(ctx context.Context, id uint64) (val lptypes.Pool, found bool)
	CalcOutAmtGivenIn(ctx sdk.Context, pool lptypes.Pool, tokenIn sdk.Coin, tokenOutDenom string, spreadFactor math.LegacyDec) (tokenOut sdk.Coin, err error)
	CalcInAmtGivenOut(ctx sdk.Context, pool lptypes.Pool, tokenOut sdk.Coin, tokenInDenom string, spreadFactor math.LegacyDec) (sdk.Coin, error)
	SwapExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, pool lptypes.Pool, tokenIn sdk.Coin, tokenOutDenom string, tokenOutMinAmount math.Int, spreadFactor math.LegacyDec) (tokenOutAmount math.Int, err error)
	SwapExactAmountOut(ctx sdk.Context, sender sdk.AccAddress, pool lptypes.Pool, tokenInDenom string, tokenInMaxAmount math.Int, tokenOut sdk.Coin, spreadFactor math.LegacyDec) (tokenInAmount math.Int, err error)
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
