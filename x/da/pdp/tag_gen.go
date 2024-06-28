package pdp

import (
	"golang.org/x/crypto/sha3"
	"math/big"

	"github.com/klauspost/reedsolomon"
)

func TagGen(blob []byte, public Public) (t *big.Int, shardSize uint64, shardCount int) {
	shardSize, shardCount, shards := ErasureCode(blob, 6)

	hash256 := sha3.New256()
	t = big.NewInt(0)

	for _, share := range shards {
		hash := hash256.Sum(share)

		t = t.Add(t, new(big.Int).SetBytes(hash))
		t = t.Mod(t, &public.Q)
	}

	return t, shardSize, shardCount
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
	shareSizeInt := length / shardCountHalf
	shardSize = uint64(shareSizeInt)

	for i := 0; i < shardCountHalf-1; i++ {
		shards[i] = make([]byte, shareSizeInt)
		copy(shards[i], blob[i*shareSizeInt:(i+1)*shareSizeInt])
	}
	i := shardCountHalf - 1
	shards[i] = make([]byte, shareSizeInt)
	copy(shards[i], blob[i*shareSizeInt:])

	for i := 0; i < shardCountHalf; i++ {
		shards[shardCountHalf+i] = make([]byte, shareSizeInt)
	}

	err = encoder.Encode(shards)
	if err != nil {
		panic(err)
	}

	return
}
