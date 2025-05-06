package merkle

import (
	"bytes"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

func Hash(data []byte) [32]byte {
	return [32]byte(poseidon.Sum(data))
}

func MerkleRoot(leaves [][32]byte) [32]byte {
	n := len(leaves)
	if n == 0 {
		return [32]byte{}
	}

	if n == 1 {
		return leaves[0]
	}

	if n&(n-1) != 0 {
		minimumPowerOfTwo := 1
		for minimumPowerOfTwo < n {
			minimumPowerOfTwo <<= 1
		}
		leaves = append(leaves, make([][32]byte, minimumPowerOfTwo-n)...)
	}

	nodes := make([][32]byte, n/2)

	for i := range nodes {
		nodes[i] = Hash(append(leaves[2*i][:], leaves[2*i+1][:]...))
	}

	return MerkleRoot(nodes)
}

func InclusionProof(leaf [32]byte, path [][33]byte, root [32]byte) bool {
	for _, p := range path {
		if p[0] == 0x00 {
			leaf = Hash(append(p[1:], leaf[:]...))
		} else {
			leaf = Hash(append(leaf[:], p[1:]...))
		}
	}
	return bytes.Equal(leaf[:], root[:])
}
