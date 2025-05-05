package types

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams creates a new Params instance.
func NewParams(
	gasPerShard uint64,
	gasChallengeUnavailability uint64,
	declarationPeriod time.Duration,
	preservationPeriod time.Duration,
	challengeResponsePeriod time.Duration,
) Params {
	return Params{
		GasPerShard:                gasPerShard,
		GasChallengeUnavailability: gasChallengeUnavailability,
		DeclarationPeriod:          declarationPeriod,
		PreservationPeriod:         preservationPeriod,
		ChallengeResponsePeriod:    challengeResponsePeriod,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		1000,
		1000000,
		time.Hour*24,
		time.Hour*24*12,
		time.Hour,
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.GasPerShard <= 0 {
		return errorsmod.Wrap(ErrInvalidPublishDataGas, "gas per shard must be positive")
	}

	if p.GasChallengeUnavailability <= 0 {
		return errorsmod.Wrap(ErrInvalidPublishDataGas, "gas challenge unavailability must be positive")
	}

	if p.DeclarationPeriod <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "declaration period must be positive")
	}

	if p.PreservationPeriod <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "preservation period must be positive")
	}

	if p.ChallengeResponsePeriod <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "challenge response period must be positive")
	}

	return nil
}
