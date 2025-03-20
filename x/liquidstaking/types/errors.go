package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/liquidstaking module sentinel errors
var (
	ErrInvalidSigner        = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidRewardPeriod  = errors.Register(ModuleName, 1101, "invalid reward period")
	ErrInvalidUnbondedDenom = errors.Register(ModuleName, 1102, "invalid unbonded denom")
)
