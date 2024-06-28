package pdp

import (
	"math/big"
)

type Public struct {
	P big.Int
	Q big.Int
	G big.Int
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
