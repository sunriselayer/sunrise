package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/swap module sentinel errors
var (
	ErrInvalidSigner              = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidRoute               = errors.Register(ModuleName, 1102, "invalid route")
	ErrInvalidAmount              = errors.Register(ModuleName, 1103, "invalid amount")
	ErrLowerThanMinOutAmount      = errors.Register(ModuleName, 1104, "lower than minimum out amount")
	ErrHigherThanMaxInAmount      = errors.Register(ModuleName, 1105, "higher than maximum in amount")
	ErrUnexpectedAmountInMismatch = errors.Register(ModuleName, 1106, "unexpected amount in mismatch")
	UnknownStrategyType           = errors.Register(ModuleName, 1107, "unknown strategy type")
)
