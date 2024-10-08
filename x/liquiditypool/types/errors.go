package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquiditypool module sentinel errors
var (
	ErrInvalidSigner            = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrPoolNotFound             = sdkerrors.Register(ModuleName, 1101, "pool not found")
	ErrInvalidBaseDenom         = sdkerrors.Register(ModuleName, 1102, "invalid base denom")
	ErrInvalidQuoteDenom        = sdkerrors.Register(ModuleName, 1103, "invalid quote denom")
	ErrInvalidTokenAmounts      = sdkerrors.Register(ModuleName, 1104, "invalid token amounts")
	ErrInvalidTickers           = sdkerrors.Register(ModuleName, 1105, "invalid tickers")
	ErrNegativeTokenAmount      = sdkerrors.Register(ModuleName, 1106, "negative token amount")
	ErrSqrtPriceToTick          = sdkerrors.Register(ModuleName, 1107, "error converting sqrt price to tick")
	ErrPriceOutOfBound          = sdkerrors.Register(ModuleName, 1108, "price out of bound")
	ErrZeroLiquidity            = sdkerrors.Register(ModuleName, 1109, "zero liquidity")
	ErrInsufficientAmountPut    = sdkerrors.Register(ModuleName, 1110, "insufficient amount of tokens were put")
	ErrInvalidFirstPosition     = sdkerrors.Register(ModuleName, 1111, "invalid first position")
	ErrPositionNotFound         = sdkerrors.Register(ModuleName, 1112, "position not found")
	ErrInsufficientLiquidity    = sdkerrors.Register(ModuleName, 1113, "insufficient liquidity")
	ErrNextTickInfoNil          = sdkerrors.Register(ModuleName, 1114, "next tick info cannot be nil")
	ErrNegativeLiquidity        = sdkerrors.Register(ModuleName, 1115, "negative liquidity")
	ErrEmptyLiquidity           = sdkerrors.Register(ModuleName, 1116, "empty liquidity")
	ErrDenomDuplication         = sdkerrors.Register(ModuleName, 1117, "in and out denom duplication")
	ErrLessThanMinAmount        = sdkerrors.Register(ModuleName, 1118, "less than minimum amount")
	ErrGreaterThanMaxAmount     = sdkerrors.Register(ModuleName, 1119, "greater than maximum amount")
	ErrUnexpectedCalcAmount     = sdkerrors.Register(ModuleName, 1120, "unexpected calculated amount")
	ErrRanOutOfTicks            = sdkerrors.Register(ModuleName, 1121, "ran out of ticks")
	ErrRanOutOfIterations       = sdkerrors.Register(ModuleName, 1122, "ran out of iterations during swap")
	ErrInvalidInDenom           = sdkerrors.Register(ModuleName, 1123, "invalid in denom")
	ErrInvalidOutDenom          = sdkerrors.Register(ModuleName, 1124, "invalid out denom")
	ErrInvalidComputedSqrtPrice = sdkerrors.Register(ModuleName, 1125, "invalid computed sqrt price")
	ErrInvalidTickIndexEncoding = sdkerrors.Register(ModuleName, 1126, "invalid tick index encoding")
	ErrOverChargeGivenIn        = sdkerrors.Register(ModuleName, 1127, "over charge swap out given in")
	ErrNegativeSqrtPrice        = sdkerrors.Register(ModuleName, 1128, "negative sqrt price")
	ErrTickIndexOutOfBoundaries = sdkerrors.Register(ModuleName, 1129, "tickIndex out of boundaries")
	ErrNotEqualSqrtPrice        = sdkerrors.Register(ModuleName, 1130, "not equal computed sqrt price")
	ErrNoSqrtPriceAfterSwap     = sdkerrors.Register(ModuleName, 1131, "no advance sqrt price after swap")
	ErrInvalidSqrtPrice         = sdkerrors.Register(ModuleName, 1132, "invalid sqrt price")
	ErrFeePositionNotFound      = sdkerrors.Register(ModuleName, 1133, "fee position not found")
	ErrNotPositionOwner         = sdkerrors.Register(ModuleName, 1134, "not a position owner")
	ErrNonPositiveLiquidity     = sdkerrors.Register(ModuleName, 1135, "non-positive liquidity")
	ErrNoPosition               = sdkerrors.Register(ModuleName, 1136, "no position found for position key")
	ErrZeroShares               = sdkerrors.Register(ModuleName, 1137, "zero shares")
	ErrAccumDoesNotExist        = sdkerrors.Register(ModuleName, 1138, "accumulator does not exist")
	ErrNegRewardAddition        = sdkerrors.Register(ModuleName, 1139, "negative reward addition")
	ErrTickNotFound             = sdkerrors.Register(ModuleName, 1140, "tick not found")
	ErrZeroActualAmountBase     = sdkerrors.Register(ModuleName, 1141, "zero actual amount base")
	ErrZeroActualAmountQuote    = sdkerrors.Register(ModuleName, 1142, "zero actual amount quote")
)
