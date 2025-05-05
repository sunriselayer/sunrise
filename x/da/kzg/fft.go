package kzg

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/fft"
)

func CalculateCoefficients(elements []fr.Element) ([]fr.Element, uint64) {
	domain := fft.NewDomain(uint64(len(elements)))
	n := int(domain.Cardinality)

	if len(elements) < n {
		elements = append(elements, make([]fr.Element, n-len(elements))...)
	}

	fft.BitReverse(elements)
	domain.FFTInverse(elements, fft.DIT)
	coeffs := elements

	return coeffs, domain.Cardinality
}

// Erasure Coding
func CalculateExtendedEvaluationPoints(coeffs []fr.Element, domainCardinality uint64) []fr.Element {
	if domainCardinality < ElementsSideLengthPerShard {
		domainCardinality = ElementsSideLengthPerShard
	}
	extDomain := fft.NewDomain(2 * domainCardinality)
	m := int(extDomain.Cardinality)

	coeffs = append(coeffs, make([]fr.Element, m-len(coeffs))...)

	fft.BitReverse(coeffs)
	extendedEvals := make([]fr.Element, m)
	copy(extendedEvals, coeffs)
	extDomain.FFT(extendedEvals, fft.DIT)

	return extendedEvals
}

func EvaluationPoints(len uint64) []fr.Element {
	domain := fft.NewDomain(len)
	xPoints := make([]fr.Element, len)
	xPoints[0].SetOne()
	for i := 1; i < int(len); i++ {
		xPoints[i].Mul(&xPoints[i-1], &domain.Generator)
	}

	return xPoints
}
