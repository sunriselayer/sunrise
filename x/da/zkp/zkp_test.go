package zkp

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/stretchr/testify/require"
)

func TestZKP(t *testing.T) {
	preImage := big.NewInt(111)

	m := native_mimc.NewMiMC()
	m.Write(preImage.Bytes())
	hash := m.Sum(nil)

	// witness definition
	assignment := Circuit{
		ShardHashes:       []frontend.Variable{preImage},
		Indices:           []frontend.Variable{0},
		ShardDoubleHashes: []frontend.Variable{hash},
		Threshold:         1,
	}

	// compiles our circuit into a R1CS
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, Hollow(&assignment))
	require.NoError(t, err)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	require.NoError(t, err)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)
	publicWitness, err := witness.Public()
	require.NoError(t, err)

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	err = groth16.Verify(proof, vk, publicWitness)
	require.NoError(t, err)
}

func TestZKP_BigSize(t *testing.T) {
	threshold := 300
	shardHashes := []frontend.Variable{}
	shardDoubleHashes := []frontend.Variable{}
	indices := []frontend.Variable{}
	for i := 0; i < threshold; i++ {
		shardHash := big.NewInt(int64(i + 10000))
		m := native_mimc.NewMiMC()
		m.Write(shardHash.Bytes())
		shardDoubleHash := m.Sum(nil)
		indices = append(indices, i)
		shardHashes = append(shardHashes, shardHash)
		shardDoubleHashes = append(shardDoubleHashes, shardDoubleHash)
	}

	// witness definition
	assignment := Circuit{
		ShardHashes:       shardHashes,
		Indices:           indices,
		ShardDoubleHashes: shardDoubleHashes,
		Threshold:         threshold,
	}

	// compiles our circuit into a R1CS
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, Hollow(&assignment))
	require.NoError(t, err)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	require.NoError(t, err)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)
	publicWitness, err := witness.Public()
	require.NoError(t, err)

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	err = groth16.Verify(proof, vk, publicWitness)
	require.NoError(t, err)
}

func TestVerifyData_BigSize(t *testing.T) {
	threshold := 100000
	shardHashes := [][]byte{}
	shardDoubleHashes := [][]byte{}
	indices := []int64{}
	for i := 0; i < threshold; i++ {
		shardHash := big.NewInt(int64(i + 10000))
		m := native_mimc.NewMiMC()
		m.Write(shardHash.Bytes())
		shardDoubleHash := m.Sum(nil)
		indices = append(indices, int64(i))
		shardHashes = append(shardHashes, shardHash.Bytes())
		shardDoubleHashes = append(shardDoubleHashes, shardDoubleHash)
	}
	err := VerifyData(indices, shardHashes, shardDoubleHashes, threshold)
	require.NoError(t, err)
}

func TestProveAndVerify(t *testing.T) {
	preImage := big.NewInt(111)

	m := native_mimc.NewMiMC()
	m.Write(preImage.Bytes())
	hash := m.Sum(nil)

	// witness definition
	assignment := Circuit{
		ShardHashes:       []frontend.Variable{preImage},
		Indices:           []frontend.Variable{0},
		ShardDoubleHashes: []frontend.Variable{hash},
		Threshold:         1,
	}

	err := ProveAndVerify(assignment)
	require.NoError(t, err)
}
