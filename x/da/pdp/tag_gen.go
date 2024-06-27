package pdp

import (
	"golang.org/x/crypto/sha3"
	"math/big"
)

func TagGen(blob []byte, public Public) (t *big.Int, shareSize uint64, shareCount int) {
	shares, shareSize, shareCount := erasureCode(blob)

	hash256 := sha3.New256()
	t = big.NewInt(0)

	for _, share := range shares {
		hash := hash256.Sum(share)

		t = t.Add(t, new(big.Int).SetBytes(hash))
		t = t.Mod(t, &public.Q)
	}

	return t, shareSize, shareCount
}

func erasureCode(blob []byte) (shares [][]byte, shareSize uint64, shareCount int) {
	return
}
