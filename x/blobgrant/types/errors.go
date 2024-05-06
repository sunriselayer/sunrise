package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/grant module sentinel errors
var (
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample                     = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrInvalidBlockHeightDuration = sdkerrors.Register(ModuleName, 1104, "sample error")
)
