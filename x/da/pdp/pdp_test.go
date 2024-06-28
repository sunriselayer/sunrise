package pdp_test

import (
	"testing"

	"github.com/klauspost/reedsolomon"
	"github.com/stretchr/testify/assert"
	"github.com/sunriselayer/sunrise/x/da/pdp"
)

func TestErasureCode(t *testing.T) {
	dataLen := 60
	data := make([]byte, dataLen)

	for i := range data {
		data[i] = byte(i & 0xff)
	}
	shardCountHalf := 6

	shardSize, shardCount, shards := pdp.ErasureCode(data, shardCountHalf)

	assert.Equal(t, dataLen/shardCountHalf, int(shardSize))
	assert.Equal(t, shardCountHalf*2, shardCount)

	concat := make([]byte, 0)

	for i := range shards {
		concat = append(concat, shards[i]...)
	}

	assert.Equal(t, data, concat[:dataLen])

	encoder, err := reedsolomon.New(shardCountHalf, shardCountHalf)
	if err != nil {
		panic(err)
	}

	valid, _ := encoder.Verify(shards)

	assert.Equal(t, true, valid)
}
