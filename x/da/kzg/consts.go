package kzg

import (
	kzg "github.com/consensys/gnark-crypto/ecc/bls12-381/kzg"
)

const (
	// Blob is extended to two times bigger.
	// To utilize Ethereum's SRS (4906 evaluation points) which supports up to 65536*2 bytes,
	// the blob size is limited to 65536 bytes.
	MaxBlobSize                = 65536
	ShardSize                  = 2048
	ElementsPerShard           = ShardSize / 32
	ElementsSideLengthPerShard = 8 // sqrt(ElementsPerShard)
)

var (
	Srs = kzg.SRS{
		Pk: kzg.ProvingKey{},
		Vk: kzg.VerifyingKey{},
	}
)
