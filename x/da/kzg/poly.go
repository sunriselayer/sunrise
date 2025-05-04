package kzg

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func InterpolateLagrange(xs, ys []fr.Element) []fr.Element {
	n := len(xs)
	poly := make([]fr.Element, n)

	var one fr.Element
	one.SetOne()

	for i := 0; i < n; i++ {
		numer := []fr.Element{one}
		denom := one

		for j := 0; j < n; j++ {
			if i == j {
				continue
			}

			// numer *= (x - xs[j])
			var negXj fr.Element
			negXj.Neg(&xs[j])
			tmp := []fr.Element{negXj}
			tmp = append(tmp, one)
			numer = polyMul(numer, tmp)

			// denom *= (xs[i] - xs[j])
			var diff fr.Element
			diff.Sub(&xs[i], &xs[j])
			denom.Mul(&denom, &diff)
		}

		// scale numer by y_i / denom
		var denomInv fr.Element
		denomInv.Inverse(&denom)
		var scalar fr.Element
		scalar.Mul(&ys[i], &denomInv)
		numer = polyScale(numer, scalar)

		poly = polyAdd(poly, numer)
	}

	return poly
}

func polyAdd(a, b []fr.Element) []fr.Element {
	n := max(len(a), len(b))
	out := make([]fr.Element, n)
	for i := 0; i < n; i++ {
		if i < len(a) {
			out[i] = a[i]
		}
		if i < len(b) {
			out[i].Add(&out[i], &b[i])
		}
	}
	return out
}

func polyMul(a, b []fr.Element) []fr.Element {
	out := make([]fr.Element, len(a)+len(b)-1)
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			var tmp fr.Element
			tmp.Mul(&a[i], &b[j])
			out[i+j].Add(&out[i+j], &tmp)
		}
	}
	return out
}

func polyScale(a []fr.Element, scalar fr.Element) []fr.Element {
	out := make([]fr.Element, len(a))
	for i := 0; i < len(a); i++ {
		out[i].Mul(&a[i], &scalar)
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
