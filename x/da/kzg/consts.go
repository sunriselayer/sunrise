package kzg

import (
	kzg "github.com/consensys/gnark-crypto/ecc/bls12-381/kzg"
)

const (
	EvaluationPointCount = 32
	DataSize             = 1024
)

var (
	Srs = kzg.SRS{
		Pk: kzg.ProvingKey{},
		Vk: kzg.VerifyingKey{},
	}
)
