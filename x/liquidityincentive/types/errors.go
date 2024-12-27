package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquidityincentive module sentinel errors
var (
	ErrInvalidSigner    = sdkerrors.Register(ModuleName, 1, "expected gov account as only signer for proposal message")
	ErrInvalidParam     = sdkerrors.Register(ModuleName, 2, "invalid param")
	ErrTotalWeightGTOne = sdkerrors.Register(ModuleName, 3, "total weight is greater than 1")
	ErrInvalidWeight    = sdkerrors.Register(ModuleName, 4, "invalid weight")
)
