package pdp_test

import (
	"math/big"
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

func TestPdp(t *testing.T) {
	dataLen := 60
	data := make([]byte, dataLen)

	for i := range data {
		data[i] = byte(i & 0xff)
	}
	shardCountHalf := 6

	p := big.NewInt(263)
	q := big.NewInt(131)
	g := big.NewInt(4)
	public := pdp.Public{
		P: *p,
		Q: *q,
		G: *g,
	}

	_, _, shards := pdp.ErasureCode(data, shardCountHalf)
	_, shardCount, tag := pdp.TagGen(data, public, shardCountHalf)

	c1 := big.NewInt(100)
	c2 := 6
	k := int64(20)
	perm := pdp.RandomPermutation(shardCount, k, c2)
	permShards := make([][]byte, c2)
	for i := range permShards {
		permShards[i] = shards[perm[i]]
	}

	proof := pdp.GenProof(public, tag, *c1, permShards)

	valid := pdp.CheckProof(public, tag, *c1, proof)

	assert.Equal(t, true, valid)
}

func TestGroupOrder(t *testing.T) {
	p := big.NewInt(263)
	q := big.NewInt(131)
	g := big.NewInt(4)

	// Calculate g^q mod p
	result := new(big.Int).Exp(g, q, p)

	// Verify that the result is 1
	assert.Equal(t, big.NewInt(1), result)
}
