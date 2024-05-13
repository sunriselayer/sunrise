package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/grant module sentinel errors
var (
	ErrInvalidSigner                    = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidGasPerLiquidity           = sdkerrors.Register(ModuleName, 1101, "invalid gas per liquidity")
	ErrInvalidExpiryDuration            = sdkerrors.Register(ModuleName, 1102, "invalid expiry duration")
	ErrInvalidGrantTokenRefillThreshold = sdkerrors.Register(ModuleName, 1103, "invalid grant token refill threshold")
	ErrInvalidBlockHeightDuration       = sdkerrors.Register(ModuleName, 1104, "invalid block height duration")
)
