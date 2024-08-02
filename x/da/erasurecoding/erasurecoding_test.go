package erasurecoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErasureCode(t *testing.T) {
	blob := []byte("erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob_erasurecode_testblob")
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
}
