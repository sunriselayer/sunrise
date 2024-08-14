package erasurecoding

import (
	"bufio"
	"bytes"

	"github.com/klauspost/reedsolomon"
)

func ErasureCode(blob []byte, shardCountHalf int) (shardSize uint64, shardCount int, shards [][]byte, err error) {
	encoder, err := reedsolomon.New(shardCountHalf, shardCountHalf)
	if err != nil {
		return
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

	extendedBlob := make([]byte, shardSize*uint64(shardCountHalf))
	copy(extendedBlob, blob)

	for i := 0; i < shardCountHalf; i++ {
		shards[i] = make([]byte, shardSizeInt)
		copy(shards[i], extendedBlob[i*shardSizeInt:(i+1)*shardSizeInt])
	}

	for i := 0; i < shardCountHalf; i++ {
		shards[shardCountHalf+i] = make([]byte, shardSizeInt)
	}

	err = encoder.Encode(shards)
	if err != nil {
		return
	}

	return
}

func ReconstructAndJoinShards(shards [][]byte, blobSize int) (blob []byte, err error) {
	shardCountHalf := len(shards) / 2
	encoder, err := reedsolomon.New(shardCountHalf, shardCountHalf)
	if err != nil {
		return nil, err
	}

	err = encoder.Reconstruct(shards)
	if err != nil {
		return nil, err
	}

	return JoinShards(shards, blobSize)
}

func JoinShards(shards [][]byte, blobSize int) (blob []byte, err error) {
	shardCountHalf := len(shards) / 2
	encoder, err := reedsolomon.New(shardCountHalf, shardCountHalf)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	bufWrite := bufio.NewWriter(&b)
	err = encoder.Join(bufWrite, shards, blobSize)
	if err != nil {
		return nil, err
	}

	err = bufWrite.Flush()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
