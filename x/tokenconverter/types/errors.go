package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/tokenconverter module sentinel errors
var (
	ErrInvalidSigner           = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrExceedSelfDelegationCap = errors.Register(ModuleName, 1101, "exceeded self delegation cap")
)
