package pdp

import (
	"crypto/rand"
	"golang.org/x/crypto/sha3"
	"io"
	"math/big"
	mrand "math/rand"
)

func GenProof(public Public, t big.Int, c1 int64, permutedShares [][]byte) Proof {
	hash256 := sha3.New256()
	tBar := big.NewInt(0)

	for _, share := range permutedShares {
		hash := hash256.Sum(share)

		tBar = tBar.Add(tBar, new(big.Int).SetBytes(hash))
		tBar = tBar.Mod(tBar, &public.Q)
	}

	tLargeBar := new(big.Int).Exp(&public.G, tBar, &public.P)

	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err)
	}
	r := new(big.Int).SetBytes(buf)

	c1Big := big.NewInt(c1)
	x := new(big.Int).Exp(&public.G, r, &public.P)

	tSub := new(big.Int).Sub(&t, tBar)
	y := r.Add(r, tSub.Mul(tSub, c1Big))
	y = y.Mod(y, &public.Q)

	return Proof{
		X:         *x,
		Y:         *y,
		TLargeBar: *tLargeBar,
	}
}

func RandomPermutation(shareCount int, k int64, c2 int) []int {
	source := mrand.NewSource(k)
	rand := mrand.New(source)

	ret := make([]int, c2)

	for i := range ret {
		ret[i] = rand.Intn(shareCount)
	}

	return ret
}
