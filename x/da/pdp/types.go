package pdp

import (
	"fmt"
	"math/big"
)

type Public struct {
	P big.Int
	Q big.Int
	G big.Int
}

func (public *Public) Validate() error {
	println(public.P.String())
	println(public.Q.String())
	println(public.G.String())
	v := new(big.Int).Exp(&public.G, &public.Q, &public.P)
	r := v.Cmp(big.NewInt(1))
	println(v.String())

	if r != 0 {
		return fmt.Errorf("It is required to satisfy g^q = 1 mod p")
	}

	return nil
}

type Proof struct {
	X         big.Int
	Y         big.Int
	TLargeBar big.Int
}

type Metadata struct {
	ShardSize  uint64
	ShardCount int
	ShardUris  []string
}
