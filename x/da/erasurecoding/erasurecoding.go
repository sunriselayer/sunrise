package erasurecoding

import (
	"bufio"
	"bytes"

	"github.com/klauspost/reedsolomon"
)

func ErasureCode(blob []byte, dataShardCount int, parityShardCount int) (shardSize uint64, shardCount int, shards [][]byte, err error) {
	encoder, err := reedsolomon.New(dataShardCount, parityShardCount)
	if err != nil {
		return
	}

	shardCount = dataShardCount + parityShardCount
	shards = make([][]byte, shardCount)

	length := len(blob)
	mod := length % dataShardCount
	if mod != 0 {
		length += dataShardCount - mod
	}
	shardSizeInt := length / dataShardCount
	shardSize = uint64(shardSizeInt)

	extendedBlob := make([]byte, shardSize*uint64(dataShardCount))
	copy(extendedBlob, blob)

	for i := 0; i < dataShardCount; i++ {
		shards[i] = make([]byte, shardSizeInt)
		copy(shards[i], extendedBlob[i*shardSizeInt:(i+1)*shardSizeInt])
	}

	for i := 0; i < parityShardCount; i++ {
		shards[dataShardCount+i] = make([]byte, shardSizeInt)
	}

	err = encoder.Encode(shards)
	if err != nil {
		return
	}

	return
}

func ReconstructAndJoinShards(shards [][]byte, dataShardCount int, blobSize int) (blob []byte, err error) {
	parityShardCount := len(shards) - dataShardCount
	encoder, err := reedsolomon.New(dataShardCount, parityShardCount)
	if err != nil {
		return nil, err
	}

	err = encoder.Reconstruct(shards)
	if err != nil {
		return nil, err
	}

	return JoinShards(shards, dataShardCount, blobSize)
}

func JoinShards(shards [][]byte, dataShardCount int, blobSize int) (blob []byte, err error) {
	parityShardCount := len(shards) - dataShardCount
	encoder, err := reedsolomon.New(dataShardCount, parityShardCount)
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
