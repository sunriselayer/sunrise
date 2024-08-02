package zkp

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/stretchr/testify/require"
)

func TestZKPDemo(t *testing.T) {
	// compiles our circuit into a R1CS
	var circuit DemoCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	require.NoError(t, err)

	// witness definition
	assignment := DemoCircuit{
		PreImage: "16130099170765464552823636852555369511329944820189892919423002775646948828469",
		Hash:     "12886436712380113721405259596386800092738845035233065858332878701083870690753",
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)
	publicWitness, err := witness.Public()
	require.NoError(t, err)

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	err = groth16.Verify(proof, vk, publicWitness)
	require.NoError(t, err)
}

func TestZKPDemoWithHasher(t *testing.T) {
	// compiles our circuit into a R1CS
	var circuit DemoCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	require.NoError(t, err)

	preImage := big.NewInt(111)

	m := native_mimc.NewMiMC()
	m.Write(preImage.Bytes())

	var hash fr.Element
	hash.SetBytes(m.Sum(nil))

	fmt.Printf("Hash:%X", hash.Bytes())

	// witness definition
	assignment := DemoCircuit{
		PreImage: preImage,
		Hash:     m.Sum(nil), // hash
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)
	publicWitness, err := witness.Public()
	require.NoError(t, err)

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	err = groth16.Verify(proof, vk, publicWitness)
	require.NoError(t, err)
}

func TestZKPDemoCommitmentCircuit(t *testing.T) {
	// witness definition - five commitments five public
	assignment := CommitmentCircuit{
		X:      []frontend.Variable{0, 1, 2, 3, 4},
		Public: []frontend.Variable{1, 2, 3, 4, 5},
	}

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
