package pdp

import (
	"math/big"
)

func CheckProof(public Public, t big.Int, k uint32, c1 int64, proof Proof) bool {
	tLarge := new(big.Int).Exp(&public.G, new(big.Int).Neg(&t), &public.P)

	tPrime := new(big.Int).Mul(tLarge, &proof.TLargeBar)
	tPrime = tPrime.Mod(tPrime, &public.P)

	gPowY := new(big.Int).Exp(&public.G, &proof.Y, &public.P)
	tBarPowC1 := new(big.Int).Exp(tPrime, big.NewInt(c1), &public.P)
	x := new(big.Int).Mul(gPowY, tBarPowC1)

	return proof.X.Cmp(x) == 0
}
