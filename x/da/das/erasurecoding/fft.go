package erasurecoding

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/fft"

	"github.com/sunriselayer/sunrise/x/da/das/consts"
)

func CalculateCoefficients(elements []fr.Element) []fr.Element {
	if len(elements) == 0 {
		return []fr.Element{}
	}
	domain := fft.NewDomain(uint64(len(elements)))
	n := int(domain.Cardinality)

	if len(elements) < n {
		elements = append(elements, make([]fr.Element, n-len(elements))...)
	}

	fft.BitReverse(elements)
	domain.FFTInverse(elements, fft.DIT)
	coeffs := elements

	return coeffs
}

// Erasure Coding
func CalculateExtendedEvaluationPoints(coeffs []fr.Element) ([]fr.Element, error) {
	domainCardinality := len(coeffs)
	if domainCardinality == 0 {
		return nil, fmt.Errorf("coefficients cannot be empty")
	}

	if domainCardinality < consts.ElementsLenPerShard {
		domainCardinality = consts.ElementsLenPerShard
	}
	extDomain := fft.NewDomain(consts.ExtensionRatio * uint64(domainCardinality))
	m := int(extDomain.Cardinality)

	// Create a copy of coefficients and pad with zeros
	extendedEvals := make([]fr.Element, m)
	copy(extendedEvals, coeffs)

	// Perform FFT on the extended domain
	fft.BitReverse(extendedEvals)
	extDomain.FFT(extendedEvals, fft.DIT)

	return extendedEvals, nil
}

func EvaluationPoints(len uint64) []fr.Element {
	if len == 0 {
		return []fr.Element{}
	}
	domain := fft.NewDomain(len)
	xPoints := make([]fr.Element, len)
	xPoints[0].SetOne()
	for i := 1; i < int(len); i++ {
		xPoints[i].Mul(&xPoints[i-1], &domain.Generator)
	}

	return xPoints
}
