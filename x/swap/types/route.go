package types

import (
	"fmt"
)

func (route *Route) Validate() error {
	switch strategy := route.Strategy.(type) {
	case *Route_Pool:
		return nil
	case *Route_Series:
		series := strategy.Series

		if len(series.Routes) == 0 {
			return fmt.Errorf("TODO")
		}

		denomIn := route.DenomIn
		for _, r := range series.Routes {
			if err := r.Validate(); err != nil {
				return err
			}

			if r.DenomIn != denomIn {
				return fmt.Errorf("TODO")
			}
			denomIn = r.DenomOut
		}
		if denomIn != route.DenomOut {
			return fmt.Errorf("TODO")
		}

		return nil
	case *Route_Parallel:
		parallel := strategy.Parallel

		if len(parallel.Routes) == 0 {
			return fmt.Errorf("TODO")
		}
		if len(parallel.Routes) != len(parallel.Weights) {
			return fmt.Errorf("TODO")
		}

		for i, r := range parallel.Routes {
			if err := r.Validate(); err != nil {
				return err
			}

			if r.DenomIn != route.DenomIn {
				return fmt.Errorf("TODO")
			}
			if r.DenomOut != route.DenomOut {
				return fmt.Errorf("TODO")
			}

			if parallel.Weights[i].IsNil() {
				return fmt.Errorf("TODO")
			}
			if !parallel.Weights[i].IsPositive() {
				return fmt.Errorf("TODO")
			}
		}

		return nil
	}

	return fmt.Errorf("TODO")

}
