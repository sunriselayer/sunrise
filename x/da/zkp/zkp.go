package zkp

import (
	"bytes"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	native_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
)

// Circuit defines our ZKP circuit
type Circuit struct {
	Indices           []frontend.Variable `gnark:",public"`  // indices
	ShardHashes       []frontend.Variable `gnark:",private"` // H(s_i) values
	ShardDoubleHashes []frontend.Variable `gnark:",public"`  // H^2(s_i) values
	Threshold         frontend.Variable   `gnark:",public"`  // threshold
}

func (circuit *Circuit) Define(api frontend.API) error {
	mimcHash, err := mimc.NewMiMC(api)
	if err != nil {
		return fmt.Errorf("failed to create MiMC hash: %v", err)
	}

	// Verify hashes of shares
	for i, j := range circuit.Indices {
		mimcHash.Write(circuit.ShardHashes[i])
		h := mimcHash.Sum()
		j, _ := j.(int)
		api.AssertIsEqual(h, circuit.ShardDoubleHashes[j])
	}

	// Verify that the number of shares used is at least the threshold
	api.AssertIsLessOrEqual(circuit.Threshold, frontend.Variable(len(circuit.ShardHashes)))

	return nil
}

func ProveAndVerify(assignment Circuit) error {
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, Hollow(&assignment))
	if err != nil {
		return fmt.Errorf("compile error: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf("setup error: %v", err)
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return fmt.Errorf("witness creation error: %v", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return fmt.Errorf("proving error: %v", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return fmt.Errorf("public witness error: %v", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		return fmt.Errorf("verification failed: %v", err)
	}

	return nil
}

func VerifyDataZKP(indices []int64, shardHashes [][]byte, shardDoubleHashes [][]byte, threshold int) error {
	indicesFrontend := []frontend.Variable{}
	shardHashesFrontend := []frontend.Variable{}
	shardDoubleHashesFrontend := []frontend.Variable{}
	for _, shardHash := range shardHashes {
		shardHashesFrontend = append(shardHashesFrontend, shardHash)
	}

	for _, shardDoubleHash := range shardDoubleHashes {
		shardDoubleHashesFrontend = append(shardDoubleHashesFrontend, shardDoubleHash)
	}

	for _, indice := range indices {
		indicesFrontend = append(indicesFrontend, indice)
	}

	return ProveAndVerify(Circuit{
		Indices:           indicesFrontend,
		ShardHashes:       shardHashesFrontend,
		ShardDoubleHashes: shardDoubleHashesFrontend,
		Threshold:         threshold,
	})
}

func VerifyData(indices []int64, shardHashes [][]byte, shardDoubleHashes [][]byte, threshold int) error {
	for i, j := range indices {
		m := native_mimc.NewMiMC()
		m.Write(shardHashes[i])
		h := m.Sum(nil)
		if !bytes.Equal(h, shardDoubleHashes[j]) {
			return fmt.Errorf("hash mismatch for (%d, %d)", i, j)
		}
	}
	return nil
}

// ValidityProofCircuit defines single hash verification ZKP circuit
type ValidityProofCircuit struct {
	ShardHash       frontend.Variable `gnark:",private"` // H(s) value
	ShardDoubleHash frontend.Variable `gnark:",public"`  // H^2(s) value
}

func (circuit *ValidityProofCircuit) Define(api frontend.API) error {
	mimcHash, err := mimc.NewMiMC(api)
	if err != nil {
		return fmt.Errorf("failed to create MiMC hash: %v", err)
	}

	// Verify hash
	mimcHash.Write(circuit.ShardHash)
	h := mimcHash.Sum()
	api.AssertIsEqual(h, circuit.ShardDoubleHash)

	return nil
}
