package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/tokenconverter module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample        = sdkerrors.Register(ModuleName, 1101, "sample error")

	ErrExceedsMaxSupply = sdkerrors.Register(ModuleName, 1111, "exceeds max supply")

	ErrInsufficientAmountOut = sdkerrors.Register(ModuleName, 1121, "insufficient amount out")
	ErrExceededAmountIn      = sdkerrors.Register(ModuleName, 1122, "exceeded amount in")

	ErrEmptyGovDenom        = sdkerrors.Register(ModuleName, 1200, "gov denom cannot be empty")
	ErrEmptyFeeDenom        = sdkerrors.Register(ModuleName, 1201, "fee denom cannot be empty")
	ErrNegativeMaxSupplyFee = sdkerrors.Register(ModuleName, 1202, "max supply fee cannot be negative")
)
