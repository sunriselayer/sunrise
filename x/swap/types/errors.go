package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/swap module sentinel errors
var (
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample                     = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidRoute               = sdkerrors.Register(ModuleName, 1102, "invalid route")
	ErrInvalidAmount              = sdkerrors.Register(ModuleName, 1103, "invalid amount")
	ErrLowerThanMinOutAmount      = sdkerrors.Register(ModuleName, 1104, "lower than minimum out amount")
	ErrHigherThanMaxInAmount      = sdkerrors.Register(ModuleName, 1105, "higher than maximum in amount")
	ErrUnexpectedAmountInMismatch = sdkerrors.Register(ModuleName, 1106, "unexpected amount in mismatch")
	UnknownStrategyType           = sdkerrors.Register(ModuleName, 1107, "unknown strategy type")
)
