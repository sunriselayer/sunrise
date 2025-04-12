package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/liquidityincentive module sentinel errors
var (
	ErrInvalidSigner       = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrTotalWeightGTOne    = errors.Register(ModuleName, 2, "total weight is greater than 1")
	ErrInvalidWeight       = errors.Register(ModuleName, 3, "invalid weight")
	ErrBribeAlreadyClaimed = errors.Register(ModuleName, 4, "bribe already claimed")
)
