package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/swap module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample        = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidRoute  = sdkerrors.Register(ModuleName, 1102, "invalid route")
	ErrInvalidAmount = sdkerrors.Register(ModuleName, 1103, "invalid amount")
)
