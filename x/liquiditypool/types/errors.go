package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquiditypool module sentinel errors
var (
	ErrInvalidSigner         = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrPoolNotFound          = sdkerrors.Register(ModuleName, 1101, "pool not found")
	ErrInvalidBaseDenom      = sdkerrors.Register(ModuleName, 1102, "invalid base denom")
	ErrInvalidQuoteDenom     = sdkerrors.Register(ModuleName, 1103, "invalid quote denom")
	ErrInvalidTokenAmounts   = sdkerrors.Register(ModuleName, 1104, "invalid token amounts")
	ErrInvalidTickers        = sdkerrors.Register(ModuleName, 1105, "invalid tickers")
	ErrNegativeTokenAmount   = sdkerrors.Register(ModuleName, 1106, "negative token amount")
	ErrSqrtPriceToTick       = sdkerrors.Register(ModuleName, 1107, "error converting sqrt price to tick")
	ErrPriceOutOfBound       = sdkerrors.Register(ModuleName, 1108, "price out of bound")
	ErrZeroLiquidity         = sdkerrors.Register(ModuleName, 1109, "zero liquidity")
	ErrInsufficientAmountPut = sdkerrors.Register(ModuleName, 1110, "insufficient amount of tokens were put")
	ErrInvalidFirstPosition  = sdkerrors.Register(ModuleName, 1111, "invalid first position")
	ErrPositionNotFound      = sdkerrors.Register(ModuleName, 1112, "position not found")
	ErrInsufficientLiquidity = sdkerrors.Register(ModuleName, 1113, "insufficient liquidity")
	ErrNextTickInfoNil       = sdkerrors.Register(ModuleName, 1114, "next tick info cannot be nil")
	ErrNegativeLiquidity     = sdkerrors.Register(ModuleName, 1115, "negative liquidity")
)
