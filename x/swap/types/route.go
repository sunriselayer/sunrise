package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (route *Route) Validate() error {
	if err := route.validateRecursive(); err != nil {
		return errorsmod.Wrapf(ErrInvalidRoute, "%s", err)
	}

	// Check if the pool is reused
	// Reuse must be prevented because it causes a problem in the calculation of the slippage
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			panic(errorsmod.Wrapf(ErrInvalidRoute, "%s", err))
		}
	}()
	route.mustNotReusePool(make(map[uint64]bool))

	return nil
}

func (route *Route) validateRecursive() error {
	activeStrategies := 0
	if route.Pool != nil {
		activeStrategies++
	}
	if route.Series != nil {
		activeStrategies++
	}
	if route.Parallel != nil {
		activeStrategies++
	}
	if activeStrategies == 0 {
		return UnknownStrategyType
	} else if activeStrategies > 1 {
		return TooManyStrategyTypes
	}
	if route.Pool != nil {
		return nil
	} else if route.Series != nil {
		series := route.Series

		if len(series.Routes) == 0 {
			return fmt.Errorf("empty series")
		}

		denomIn := route.DenomIn
		for _, r := range series.Routes {
			if err := r.validateRecursive(); err != nil {
				return err
			}

			if r.DenomIn != denomIn {
				return fmt.Errorf("invalid denom in: %s", r)
			}
			denomIn = r.DenomOut
		}
		denomOut := denomIn
		if denomOut != route.DenomOut {
			return fmt.Errorf("denom out mismatch: %s, %s", denomOut, route.DenomOut)
		}

		return nil
	} else {
		parallel := route.Parallel

		if len(parallel.Routes) == 0 {
			return fmt.Errorf("empty parallel")
		}
		if len(parallel.Routes) != len(parallel.Weights) {
			return fmt.Errorf("mismatched length of parallel routes and weights")
		}

		for i, r := range parallel.Routes {
			if err := r.validateRecursive(); err != nil {
				return err
			}

			if r.DenomIn != route.DenomIn {
				return fmt.Errorf("invalid denom in: %s", r)
			}
			if r.DenomOut != route.DenomOut {
				return fmt.Errorf("invalid denom out: %s", r)
			}

			if parallel.Weights[i].IsNil() {
				return fmt.Errorf("nil weight")
			}
			if !parallel.Weights[i].IsPositive() {
				return fmt.Errorf("non-positive weight: %s", parallel.Weights[i])
			}
		}
		return nil
	}
}

func (route *Route) mustNotReusePool(poolIds map[uint64]bool) {
	if route.Pool != nil {
		poolId := route.Pool.PoolId
		if poolIds[poolId] {
			panic(fmt.Sprintf("reused pool: %d", poolId))
		}
		poolIds[poolId] = true
	} else if route.Series != nil {
		series := route.Series

		for _, r := range series.Routes {
			r.mustNotReusePool(poolIds)
		}
	} else {
		parallel := route.Parallel

		for _, r := range parallel.Routes {
			r.mustNotReusePool(poolIds)
		}
	}
}

func (route *Route) InspectRoute(
	amountExact math.Int,
	inspectRoutePool func(
		denomIn string,
		denomOut string,
		pool RoutePool,
		amountExact math.Int,
	) (amountResult math.Int, err error),
	generateResult func(
		denomIn string,
		denomOut string,
		amountExact math.Int,
		amountResult math.Int,
	) (tokenIn sdk.Coin, tokenOut sdk.Coin),
	reverse bool,
) (amountResult math.Int, routeResult RouteResult, err error) {
	if route.Pool != nil {
		amountResult, err := inspectRoutePool(
			route.DenomIn,
			route.DenomOut,
			*route.Pool,
			amountExact,
		)
		if err != nil {
			return math.Int{}, RouteResult{}, err
		}

		tokenIn, tokenOut := generateResult(route.DenomIn, route.DenomOut, amountExact, amountResult)

		return amountResult, RouteResult{
			TokenIn:  tokenIn,
			TokenOut: tokenOut,
			Strategy: &RouteResult_Pool{
				Pool: &RouteResultPool{
					PoolId: route.Pool.PoolId,
				},
			},
		}, nil

	} else if route.Series != nil {
		amountExactBuffer := amountExact
		results := make([]RouteResult, len(route.Series.Routes))
		for i := range route.Series.Routes {
			var r *Route
			if !reverse {
				r = &route.Series.Routes[i]
			} else {
				r = &route.Series.Routes[len(route.Series.Routes)-1-i]
			}
			amountResultBuffer, routeResultBuffer, err := r.InspectRoute(amountExactBuffer, inspectRoutePool, generateResult, reverse)
			if err != nil {
				return math.Int{}, RouteResult{}, err
			}
			results[i] = routeResultBuffer

			amountExactBuffer = amountResultBuffer
		}
		amountResult = amountExactBuffer

		tokenIn, tokenOut := generateResult(route.DenomIn, route.DenomOut, amountExact, amountResult)

		return amountResult, RouteResult{
			TokenIn:  tokenIn,
			TokenOut: tokenOut,
			Strategy: &RouteResult_Series{
				Series: &RouteResultSeries{
					RouteResults: results,
				},
			},
		}, nil

	} else {
		// Calculate the sum of the weights
		weightSum := math.LegacyZeroDec()
		for _, w := range route.Parallel.Weights {
			weightSum.AddMut(w)
		}

		// Calculate the amount of input for each route
		amountsExact := make([]math.Int, len(route.Parallel.Routes))
		amountsExactSum := math.ZeroInt()
		length := len(route.Parallel.Weights)

		for i, w := range route.Parallel.Weights[:length-1] {
			amountsExact[i] = w.MulInt(amountExact).Quo(weightSum).TruncateInt()
		}
		// For avoiding rounding errors
		amountsExact[length-1] = amountExact.Sub(amountsExactSum)

		// Preparations for the results
		amountResult = math.ZeroInt()
		results := make([]RouteResult, len(route.Parallel.Routes))

		// Execute the inspections
		for i, r := range route.Parallel.Routes {
			amountResultBuffer, routeResultBuffer, err := r.InspectRoute(amountsExact[i], inspectRoutePool, generateResult, reverse)
			if err != nil {
				return math.Int{}, RouteResult{}, err
			}
			amountResult = amountResult.Add(amountResultBuffer)
			results[i] = routeResultBuffer
		}

		tokenIn, tokenOut := generateResult(route.DenomIn, route.DenomOut, amountExact, amountResult)

		return amountResult, RouteResult{
			TokenIn:  tokenIn,
			TokenOut: tokenOut,
			Strategy: &RouteResult_Parallel{
				Parallel: &RouteResultParallel{
					RouteResults: results,
				},
			},
		}, nil
	}
}
