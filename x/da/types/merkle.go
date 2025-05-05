package types

import (
	"bytes"
)

func InclusionProof(leaf [32]byte, path [][33]byte, root []byte, hash func([]byte) [32]byte) bool {
	for _, p := range path {
		if p[0] == 0x00 {
			leaf = hash(append(p[1:], leaf[:]...))
		} else {
			leaf = hash(append(leaf[:], p[1:]...))
		}
	}
	return bytes.Equal(leaf[:], root)
}
