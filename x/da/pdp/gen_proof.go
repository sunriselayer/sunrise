package pdp

import (
	"crypto/rand"
	"golang.org/x/crypto/sha3"
	"io"
	"math/big"
	mrand "math/rand"
)

func GenProof(public Public, t big.Int, c1 big.Int, permutedShards [][]byte) Proof {
	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}
	r := new(big.Int).SetBytes(buf)

	return GenProofWithRandom(public, t, c1, permutedShards, *r)
}

func GenProofWithRandom(public Public, t big.Int, c1 big.Int, permutedShards [][]byte, r big.Int) Proof {
	if err := public.Validate(); err != nil {
		panic(err)
	}

	hash256 := sha3.New256()
	tBar := big.NewInt(0)

	for _, shard := range permutedShards {
		hash := hash256.Sum(shard)

		tBar = tBar.Add(tBar, new(big.Int).SetBytes(hash))
		tBar = tBar.Mod(tBar, &public.Q)
	}

	x := new(big.Int).Exp(&public.G, &r, &public.P)

	tSub := t.Sub(&t, tBar)
	y := r.Add(&r, tSub.Mul(tSub, &c1))
	y = y.Mod(y, &public.Q)

	return Proof{
		X:    *x,
		Y:    *y,
		TBar: *tBar,
	}
}

func RandomPermutation(shardCount int, k int64, c2 int) []int {
	source := mrand.NewSource(k)
	rand := mrand.New(source)

	ret := make([]int, c2)

	for i := range ret {
		ret[i] = rand.Intn(shardCount)
	}

	return ret
}
