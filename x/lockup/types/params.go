package types

import (
	"time"
)

// NewParams creates a new Params instance.
func NewParams(minLockupDuration time.Duration) Params {
	return Params{
		MinLockupDuration: minLockupDuration,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(time.Hour * 24 * 30)
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
