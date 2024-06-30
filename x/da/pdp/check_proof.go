package pdp

import (
	"math/big"
)

func CheckProof(public Public, t big.Int, c1 big.Int, proof Proof) bool {
	if err := public.Validate(); err != nil {
		panic(err)
	}

	tSub := t.Sub(&proof.TBar, &t)
	tPrime := tSub.Exp(&public.G, tSub, &public.P)

	gPowY := new(big.Int).Exp(&public.G, &proof.Y, &public.P)
	tPrimePowC1 := tPrime.Exp(tPrime, &c1, &public.P)
	x := gPowY.Mul(gPowY, tPrimePowC1)
	x = x.Mod(x, &public.P)

	return proof.X.Cmp(x) == 0
}
