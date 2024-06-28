package pdp

import (
	"golang.org/x/crypto/sha3"
	"math/big"

	"github.com/klauspost/reedsolomon"
)

func TagGen(blob []byte, public Public, shardCountHalf int) (shardSize uint64, shardCount int, tag big.Int) {
	if err := public.Validate(); err != nil {
		panic(err)
	}

	shardSize, shardCount, shards := ErasureCode(blob, shardCountHalf)

	hash256 := sha3.New256()
	t := big.NewInt(0)

	for _, shard := range shards {
		hash := hash256.Sum(shard)

		t = t.Add(t, new(big.Int).SetBytes(hash))
		t = t.Mod(t, &public.Q)
	}
	tag = *t

	return shardSize, shardCount, tag
}

func ErasureCode(blob []byte, shardCountHalf int) (shardSize uint64, shardCount int, shards [][]byte) {
	encoder, err := reedsolomon.New(shardCountHalf, shardCountHalf)
	if err != nil {
		panic(err)
	}

	shardCount = shardCountHalf * 2
	shards = make([][]byte, shardCount)

	length := len(blob)
	mod := length % shardCountHalf
	if mod != 0 {
		length += shardCountHalf - mod
	}
	shardSizeInt := length / shardCountHalf
	shardSize = uint64(shardSizeInt)

	for i := 0; i < shardCountHalf-1; i++ {
		shards[i] = make([]byte, shardSizeInt)
		copy(shards[i], blob[i*shardSizeInt:(i+1)*shardSizeInt])
	}
	i := shardCountHalf - 1
	shards[i] = make([]byte, shardSizeInt)
	copy(shards[i], blob[i*shardSizeInt:])

	for i := 0; i < shardCountHalf; i++ {
		shards[shardCountHalf+i] = make([]byte, shardSizeInt)
	}

	err = encoder.Encode(shards)
	if err != nil {
		panic(err)
	}

	return
}
