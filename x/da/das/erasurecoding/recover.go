package erasurecoding

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func RecoverCoefficients(extendedElements map[uint64][]fr.Element, domainCardinality uint64) ([]fr.Element, error) {
	// Check if domainCardinality is a power of 2
	if domainCardinality&(domainCardinality-1) != 0 {
		return nil, fmt.Errorf("domainCardinality must be a power of 2")
	}

	// Check if extendedElements length matches domainCardinality
	if uint64(len(extendedElements)) != domainCardinality {
		return nil, fmt.Errorf("extendedElements length (%d) must match domainCardinality (%d)", len(extendedElements), domainCardinality)
	}

	// Get evaluation points
	xPoints := EvaluationPoints(domainCardinality)

	// Collect evaluation points and their values
	points := make([]fr.Element, 0, domainCardinality)
	values := make([]fr.Element, 0, domainCardinality)
	for i := uint64(0); i < domainCardinality; i++ {
		if elems, ok := extendedElements[i]; ok {
			if len(elems) > 0 {
				points = append(points, xPoints[i])
				values = append(values, elems[0])
			}
		}
	}

	// If we don't have enough points, we can't recover the coefficients
	if uint64(len(points)) < domainCardinality {
		return nil, fmt.Errorf("not enough evaluation points to recover coefficients")
	}

	// Initialize coefficients array
	coeffs := make([]fr.Element, domainCardinality)

	// For each coefficient
	for i := uint64(0); i < domainCardinality; i++ {
		var sum fr.Element
		// For each evaluation point
		for j := uint64(0); j < domainCardinality; j++ {
			// Calculate Lagrange basis polynomial
			var basis fr.Element
			basis.SetOne()
			for k := uint64(0); k < domainCardinality; k++ {
				if k != j {
					var diff fr.Element
					diff.Sub(&points[i], &points[k])
					var denom fr.Element
					denom.Sub(&points[j], &points[k])
					var term fr.Element
					term.Div(&diff, &denom)
					basis.Mul(&basis, &term)
				}
			}
			// Multiply by the value and add to sum
			var term fr.Element
			term.Mul(&basis, &values[j])
			sum.Add(&sum, &term)
		}
		coeffs[i] = sum
	}

	return coeffs, nil
}

func RecoverElements(coeffs []fr.Element, elemsLen uint64) []fr.Element {
	// Get evaluation points for the original domain
	xPoints := EvaluationPoints(elemsLen)

	// Initialize elements array
	elements := make([]fr.Element, elemsLen)

	// For each evaluation point
	for i := uint64(0); i < elemsLen; i++ {
		var sum fr.Element
		// For each coefficient
		for j := uint64(0); j < uint64(len(coeffs)); j++ {
			// Calculate x^j
			var power fr.Element
			power.SetOne()
			for k := uint64(0); k < j; k++ {
				power.Mul(&power, &xPoints[i])
			}
			// Multiply by coefficient and add to sum
			var term fr.Element
			term.Mul(&power, &coeffs[j])
			sum.Add(&sum, &term)
		}
		elements[i] = sum
	}

	return elements
}
