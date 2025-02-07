package types

// DONTCOVER

import (
	"cosmossdk.io/errors/v2"
)

// x/da module sentinel errors
var (
	ErrInvalidSigner            = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrCanNotOpenChallenge      = errors.Register(ModuleName, 1101, "challenge is only open for vote extension passed data")
	ErrChallengePeriodIsOver    = errors.Register(ModuleName, 1102, "challenge period is over")
	ErrDataNotInChallenge       = errors.Register(ModuleName, 1103, "data is not in challenge")
	ErrProofPeriodIsOver        = errors.Register(ModuleName, 1104, "proof period is over")
	ErrProofIndicesOverflow     = errors.Register(ModuleName, 1105, "proof indices overflow")
	ErrIndicesAndProofsMismatch = errors.Register(ModuleName, 1106, "indices and proofs count mismatch")
	ErrParityShardCountGTETotal = errors.Register(ModuleName, 1107, "parity shard count is greater than total")
	ErrInvalidIndices           = errors.Register(ModuleName, 1108, "invalid indices")
	ErrDataAlreadyExist         = errors.Register(ModuleName, 1109, "published data already exist")
)
