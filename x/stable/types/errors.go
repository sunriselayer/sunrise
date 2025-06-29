package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/stable module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	
)
