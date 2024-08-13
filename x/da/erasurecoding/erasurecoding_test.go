package erasurecoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErasureCode(t *testing.T) {
	blob := []byte("erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob")
	require.Len(t, blob, 125)

	shardCountHalf := int(3)
	shardSize, shardCount, shards := ErasureCode(blob, shardCountHalf)
	require.Equal(t, shardCount, int(shardCountHalf*2))
	require.Equal(t, shardSize, uint64(42))
	require.Len(t, shards, 6)
	require.Len(t, shards[0], int(shardSize))
	require.Len(t, shards[1], int(shardSize))
	require.Len(t, shards[2], int(shardSize))
	require.Len(t, shards[3], int(shardSize))
	require.Len(t, shards[4], int(shardSize))
	require.Len(t, shards[5], int(shardSize))

	decoded, err := JoinShards(shards, len(blob))
	require.NoError(t, err)
	require.Equal(t, string(decoded), string(blob))
}

func TestErasureCodeRecoveryThreshold(t *testing.T) {
	blob := []byte("erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob")
	require.Len(t, blob, 125)

	shardCountHalf := int(3)
	shardSize, shardCount, shards := ErasureCode(blob, shardCountHalf)
	require.Equal(t, shardCount, int(shardCountHalf*2))
	require.Equal(t, shardSize, uint64(42))
	require.Len(t, shards, 6)
	require.Len(t, shards[0], int(shardSize))
	require.Len(t, shards[1], int(shardSize))
	require.Len(t, shards[2], int(shardSize))
	require.Len(t, shards[3], int(shardSize))
	require.Len(t, shards[4], int(shardSize))
	require.Len(t, shards[5], int(shardSize))

	// One data shard's broken
	brokenShards := append([][]byte{}, shards...)
	brokenShards[0] = nil
	decoded, err := ReconstructAndJoinShards(brokenShards, len(blob))
	require.NoError(t, err)
	require.Equal(t, string(decoded), string(blob))

	// Two data shard's broken
	brokenShards = append([][]byte{}, shards...)
	brokenShards[0] = nil
	brokenShards[1] = nil
	decoded, err = ReconstructAndJoinShards(brokenShards, len(blob))
	require.NoError(t, err)
	require.Equal(t, string(decoded), string(blob))

	// Three data shard's broken
	brokenShards = append([][]byte{}, shards...)
	brokenShards[0] = nil
	brokenShards[1] = nil
	brokenShards[2] = nil
	decoded, err = ReconstructAndJoinShards(brokenShards, len(blob))
	require.NoError(t, err)
	require.Equal(t, string(decoded), string(blob))

	// Three data shard + one parity shard's broken
	brokenShards = append([][]byte{}, shards...)
	brokenShards[0] = nil
	brokenShards[1] = nil
	brokenShards[2] = nil
	brokenShards[3] = nil
	_, err = ReconstructAndJoinShards(brokenShards, len(blob))
	require.Error(t, err)

	// Three parity shard's broken
	brokenShards = append([][]byte{}, shards...)
	brokenShards[3] = nil
	brokenShards[4] = nil
	brokenShards[5] = nil
	decoded, err = ReconstructAndJoinShards(brokenShards, len(blob))
	require.NoError(t, err)
	require.Equal(t, string(decoded), string(blob))

	// Three parity shard + one data shard
	brokenShards = append([][]byte{}, shards...)
	brokenShards[2] = nil
	brokenShards[3] = nil
	brokenShards[4] = nil
	brokenShards[5] = nil
	_, err = ReconstructAndJoinShards(brokenShards, len(blob))
	require.Error(t, err)
}
