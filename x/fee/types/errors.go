package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/fee module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidPool   = errors.Register(ModuleName, 1101, "invalid pool")
	ErrPoolNotFound  = errors.Register(ModuleName, 1102, "pool not found")
)
