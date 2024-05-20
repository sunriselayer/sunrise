package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquiditypool module sentinel errors
var (
	ErrInvalidSigner       = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrPoolNotFound        = sdkerrors.Register(ModuleName, 1101, "pool not found")
	ErrInvalidBaseDenom    = sdkerrors.Register(ModuleName, 1102, "invalid base denom")
	ErrInvalidQuoteDenom   = sdkerrors.Register(ModuleName, 1103, "invalid quote denom")
	ErrInvalidTokenAmounts = sdkerrors.Register(ModuleName, 1104, "invalid token amounts")
)
