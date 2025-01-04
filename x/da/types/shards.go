package types

import (
	"math/rand/v2"

	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidatorSeed(valAddr sdk.ValAddress) uint64 {
	m := native_mimc.NewMiMC()
	m.Write(valAddr)
	hash := m.Sum(nil)
	seed := sdk.BigEndianToUint64(hash)
	return seed
}

func ShardIndicesForValidator(valAddr sdk.ValAddress, threshold, n int64) []int64 {
	seed := ValidatorSeed(valAddr)
	return GetRandomIndicesFromSeed(n, threshold, seed, 1024)
}

func GetRandomIndicesFromSeed(n, threshold int64, seed1, seed2 uint64) []int64 {
	if threshold > n {
		threshold = n
	}
	arr := []int64{}
	for i := int64(0); i < n; i++ {
		arr = append(arr, i)
	}

	s3 := rand.NewPCG(seed1, seed2)
	r3 := rand.New(s3)

	r3.Shuffle(int(n), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})

	// Return the first threshold elements from the shuffled array
	return arr[:threshold]
}
