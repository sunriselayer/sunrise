package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/stable module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidDenom  = errors.Register(ModuleName, 1101, "invalid denom")
	ErrInvalidAmount = errors.Register(ModuleName, 1102, "invalid amount")
	ErrDenomNotFound = errors.Register(ModuleName, 1103, "denom metadata not found")
)
