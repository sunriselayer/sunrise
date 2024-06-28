package pdp

import (
	"math/big"
)

func CheckProof(public Public, t big.Int, c1 big.Int, proof Proof) bool {
	if err := public.Validate(); err != nil {
		panic(err)
	}

	tLarge := new(big.Int).Exp(&public.G, new(big.Int).Neg(&t), &public.P)

	tPrime := new(big.Int).Mul(tLarge, &proof.TLargeBar)
	tPrime = tPrime.Mod(tPrime, &public.P)

	gPowY := new(big.Int).Exp(&public.G, &proof.Y, &public.P)
	tBarPowC1 := new(big.Int).Exp(tPrime, &c1, &public.P)
	x := new(big.Int).Mul(gPowY, tBarPowC1)

	return proof.X.Cmp(x) == 0
}
