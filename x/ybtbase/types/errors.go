package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/ybtbase module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidRequest = errors.Register(ModuleName, 1101, "invalid request")
	ErrTokenAlreadyExists = errors.Register(ModuleName, 1102, "token already exists")
	ErrTokenNotFound = errors.Register(ModuleName, 1103, "token not found")
	ErrUnauthorized = errors.Register(ModuleName, 1104, "unauthorized")
	ErrInsufficientBalance = errors.Register(ModuleName, 1105, "insufficient balance")
	ErrNoPermission = errors.Register(ModuleName, 1106, "no permission to claim yield")
)
