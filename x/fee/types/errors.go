package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/fee module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample        = sdkerrors.Register(ModuleName, 1101, "sample error")

	ErrEmptyFeeDenom    = sdkerrors.Register(ModuleName, 1200, "fee denom cannot be empty")
	ErrInvalidBurnRatio = sdkerrors.Register(ModuleName, 1201, "burn ratio must be positive and less than 1")
	ErrEmptyBypassDenom = sdkerrors.Register(ModuleName, 1202, "bypass denom cannot be empty")
)
