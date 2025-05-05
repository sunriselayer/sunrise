package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/da module sentinel errors
var (
	ErrInvalidSigner            = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrDeclarationAlreadyExists = errors.Register(ModuleName, 1101, "declaration already exists")
)
