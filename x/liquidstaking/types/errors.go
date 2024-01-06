package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquidstaking module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")

	ErrNoValidatorFound           = sdkerrors.Register(ModuleName, 2, "validator does not exist")
	ErrNoDelegatorForAddress      = sdkerrors.Register(ModuleName, 3, "delegator does not contain delegation")
	ErrInvalidDenom               = sdkerrors.Register(ModuleName, 4, "invalid denom")
	ErrNotEnoughDelegationShares  = sdkerrors.Register(ModuleName, 5, "not enough delegation shares")
	ErrRedelegationsNotCompleted  = sdkerrors.Register(ModuleName, 6, "active redelegations cannot be transferred")
	ErrUntransferableShares       = sdkerrors.Register(ModuleName, 7, "shares cannot be transferred")
	ErrSelfDelegationBelowMinimum = sdkerrors.Register(ModuleName, 8, "validator's self delegation must be greater than their minimum self delegation")
)
