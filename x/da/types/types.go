package types

import (
	"encoding/binary"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/sunriselayer/sunrise/x/da/das/consts"
)

func CalculateShardCountPerValidator(shardCount uint32, validatorCount uint32) uint32 {
	return uint32(math.Min(
		float64(shardCount),
		math.Ceil(float64(ValidatorReplicationFactor)*float64(shardCount)/float64(validatorCount)),
	))
}

func CorrespondingShardIndices(
	shardsMerkleRoot []byte,
	addr sdk.ValAddress,
	shardCount uint32,
	shardCountPerValidator uint32,
) (map[uint32]bool, error) {
	indices := make(map[uint32]bool)

	i := uint32(0)
	j := uint32(0)
	for {
		hasher, err := poseidon.New(16)
		if err != nil {
			return nil, err
		}
		hasher.Write(shardsMerkleRoot)
		hasher.Write(addr.Bytes())
		hasher.Write(uint32ToBytes(i))
		hash := hasher.Sum(nil)

		hashInt := new(big.Int).SetBytes(hash)
		hashInt.Mod(hashInt, big.NewInt(int64(shardCount)))

		index := uint32(hashInt.Uint64())
		if _, ok := indices[index]; !ok {
			indices[index] = true
			j++

			if j == shardCountPerValidator {
				break
			}
		}

		i++
	}

	return indices, nil
}

func uint32ToBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

func CorrespondingEvaluationPointIndex(
	shardsMerkleRoot []byte,
	addr sdk.ValAddress,
) (uint32, error) {
	hasher, err := poseidon.New(16)
	if err != nil {
		return 0, err
	}
	hasher.Write(shardsMerkleRoot)
	hasher.Write(addr.Bytes())
	hash := hasher.Sum(nil)

	hashInt := new(big.Int).SetBytes(hash)
	hashInt.Mod(hashInt, big.NewInt(int64(consts.ElementsLenPerShard)))

	return uint32(hashInt.Uint64()), nil
}
