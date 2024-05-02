package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquiditypool module sentinel errors
var (
	ErrInvalidSigner                = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInsufficientFootprintForTwap = sdkerrors.Register(ModuleName, 1101, "price footprint must be longer than 1 to calculate twap")
)
