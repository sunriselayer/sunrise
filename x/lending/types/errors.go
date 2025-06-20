package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/lending module sentinel errors
var (
	ErrInvalidSigner               = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidLtvRatio             = errors.Register(ModuleName, 1101, "invalid ltv ratio")
	ErrInvalidLiquidationThreshold = errors.Register(ModuleName, 1102, "invalid liquidation threshold")
	ErrInvalidInterestRate         = errors.Register(ModuleName, 1103, "invalid interest rate")
	ErrMarketNotFound              = errors.Register(ModuleName, 1104, "market not found")
	ErrInsufficientBalance         = errors.Register(ModuleName, 1105, "insufficient balance")
	ErrInvalidAmount               = errors.Register(ModuleName, 1106, "invalid amount")
	ErrUserPositionNotFound        = errors.Register(ModuleName, 1107, "user position not found")
	ErrBorrowNotFound              = errors.Register(ModuleName, 1108, "borrow not found")
	ErrInsufficientCollateral      = errors.Register(ModuleName, 1109, "insufficient collateral")
	ErrUndercollateralized         = errors.Register(ModuleName, 1110, "position is undercollateralized")
)
