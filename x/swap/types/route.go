package types

import (
	errorsmod "cosmossdk.io/errors"
)

func (route *Route) Validate() error {
	switch strategy := route.Strategy.(type) {
	case *Route_Pool:
		return nil
	case *Route_Series:
		series := strategy.Series

		if len(series.Routes) == 0 {
			return errorsmod.Wrapf(ErrInvalidRoute, "empty series")
		}

		denomIn := route.DenomIn
		for _, r := range series.Routes {
			if err := r.Validate(); err != nil {
				return err
			}

			if r.DenomIn != denomIn {
				return errorsmod.Wrapf(ErrInvalidRoute, "invalid denom in: %s", r)
			}
			denomIn = r.DenomOut
		}
		denomOut := denomIn
		if denomOut != route.DenomOut {
			return errorsmod.Wrapf(ErrInvalidRoute, "denom out mismatch: %s, %s", denomOut, route.DenomOut)
		}

		return nil
	case *Route_Parallel:
		parallel := strategy.Parallel

		if len(parallel.Routes) == 0 {
			return errorsmod.Wrapf(ErrInvalidRoute, "empty parallel")
		}
		if len(parallel.Routes) != len(parallel.Weights) {
			return errorsmod.Wrapf(ErrInvalidRoute, "mismatched length of parallel routes and weights")
		}

		for i, r := range parallel.Routes {
			if err := r.Validate(); err != nil {
				return err
			}

			if r.DenomIn != route.DenomIn {
				return errorsmod.Wrapf(ErrInvalidRoute, "invalid denom in: %s", r)
			}
			if r.DenomOut != route.DenomOut {
				return errorsmod.Wrapf(ErrInvalidRoute, "invalid denom out: %s", r)
			}

			if parallel.Weights[i].IsNil() {
				return errorsmod.Wrapf(ErrInvalidRoute, "nil weight")
			}
			if !parallel.Weights[i].IsPositive() {
				return errorsmod.Wrapf(ErrInvalidRoute, "non-positive weight: %s", parallel.Weights[i])
			}
		}

		return nil
	}

	return errorsmod.Wrapf(ErrInvalidRoute, "unknown strategy: %s", route.Strategy)
}
