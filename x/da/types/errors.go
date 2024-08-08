package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/da module sentinel errors
var (
	ErrInvalidSigner         = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrCanNotOpenChallenge   = sdkerrors.Register(ModuleName, 1101, "challenge is only open for vote extension passed data")
	ErrChallengePeriodIsOver = sdkerrors.Register(ModuleName, 1102, "challenge period is over")
	ErrDataNotInChallenge    = sdkerrors.Register(ModuleName, 1103, "data is not in challenge")
	ErrProofPeriodIsOver     = sdkerrors.Register(ModuleName, 1104, "proof period is over")
)
