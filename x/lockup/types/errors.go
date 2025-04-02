package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/lockup module sentinel errors
var (
	ErrInvalidSigner              = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidTimeRange           = errors.Register(ModuleName, 1101, "invalid time range")
	ErrLockupAccountAlreadyExists = errors.Register(ModuleName, 1102, "lockup account already exists")
	ErrInsufficientUnlockedFunds  = errors.Register(ModuleName, 1103, "insufficient unlocked funds")
)
