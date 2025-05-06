package consts

import (
	kzg "github.com/consensys/gnark-crypto/ecc/bls12-381/kzg"
)

const (
	SrsLen              = 4096
	ElementSize         = 32
	MaxRowSize          = SrsLen * ElementSize
	ShardSize           = 2048
	ElementsLenPerShard = ShardSize / ElementSize
)

var (
	Srs = kzg.SRS{
		Pk: kzg.ProvingKey{},
		Vk: kzg.VerifyingKey{},
	}
)
