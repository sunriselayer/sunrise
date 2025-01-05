package types

// NewParams creates a new Params instance.
func NewParams(bondDenom string, feeDenom string) Params {
	return Params{
		BondDenom: bondDenom,
		FeeDenom:  feeDenom,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		"stake",
		"fee",
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
