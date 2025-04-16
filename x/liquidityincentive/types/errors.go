package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/liquidityincentive module sentinel errors
var (
	ErrInvalidSigner          = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrTotalWeightGTOne       = errors.Register(ModuleName, 2, "total weight is greater than 1")
	ErrInvalidWeight          = errors.Register(ModuleName, 3, "invalid weight")
	ErrNoBribesToClaim        = errors.Register(ModuleName, 4, "no bribes to claim")
	ErrInvalidBribe           = errors.Register(ModuleName, 5, "invalid bribe")
	ErrBribeAlreadyExists     = errors.Register(ModuleName, 6, "bribe already exists")
	ErrBribeNotFound          = errors.Register(ModuleName, 7, "bribe not found")
	ErrInvalidClaimAmount     = errors.Register(ModuleName, 8, "invalid claim amount")
	ErrInsufficientBribeFunds = errors.Register(ModuleName, 9, "insufficient bribe funds")
)
