package wasmclient

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

// CurveID represents the elliptic curve used for the proof
type CurveID uint8

const (
	CurveBN254 CurveID = iota
	CurveBLS12381
)

// VerifyGroth16WithCurve verifies a Groth16 proof with a specific curve
func VerifyGroth16WithCurve(proofBytes, verifyingKeyBytes, publicWitnessBytes []byte, curveID CurveID) error {
	// Validate inputs
	if len(proofBytes) == 0 {
		return errors.New("proof bytes cannot be empty")
	}
	if len(verifyingKeyBytes) == 0 {
		return errors.New("verifying key bytes cannot be empty")
	}
	if len(publicWitnessBytes) == 0 {
		return errors.New("public witness bytes cannot be empty")
	}

	// Verify based on curve
	switch curveID {
	case CurveBN254:
		return verifyBN254(proofBytes, verifyingKeyBytes, publicWitnessBytes)
	case CurveBLS12381:
		return verifyBLS12381(proofBytes, verifyingKeyBytes, publicWitnessBytes)
	default:
		return fmt.Errorf("unsupported curve ID: %d", curveID)
	}
}

// verifyBN254 verifies a proof using BN254 curve
func verifyBN254(proofBytes, vkBytes, publicWitnessBytes []byte) error {
	// Deserialize proof
	proof := groth16.NewProof(ecc.BN254)
	if _, err := proof.ReadFrom(bytes.NewReader(proofBytes)); err != nil {
		return fmt.Errorf("failed to deserialize proof: %w", err)
	}

	// Deserialize verifying key
	vk := groth16.NewVerifyingKey(ecc.BN254)
	if _, err := vk.ReadFrom(bytes.NewReader(vkBytes)); err != nil {
		return fmt.Errorf("failed to deserialize verifying key: %w", err)
	}

	// Deserialize public witness
	publicWitness, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return fmt.Errorf("failed to create witness: %w", err)
	}
	if _, err := publicWitness.ReadFrom(bytes.NewReader(publicWitnessBytes)); err != nil {
		return fmt.Errorf("failed to deserialize public witness: %w", err)
	}

	// Verify the proof
	return groth16.Verify(proof, vk, publicWitness)
}

// verifyBLS12381 verifies a proof using BLS12-381 curve
func verifyBLS12381(proofBytes, vkBytes, publicWitnessBytes []byte) error {
	// Deserialize proof
	proof := groth16.NewProof(ecc.BLS12_381)
	if _, err := proof.ReadFrom(bytes.NewReader(proofBytes)); err != nil {
		return fmt.Errorf("failed to deserialize proof: %w", err)
	}

	// Deserialize verifying key
	vk := groth16.NewVerifyingKey(ecc.BLS12_381)
	if _, err := vk.ReadFrom(bytes.NewReader(vkBytes)); err != nil {
		return fmt.Errorf("failed to deserialize verifying key: %w", err)
	}

	// Deserialize public witness
	publicWitness, err := witness.New(ecc.BLS12_381.ScalarField())
	if err != nil {
		return fmt.Errorf("failed to create witness: %w", err)
	}
	if _, err := publicWitness.ReadFrom(bytes.NewReader(publicWitnessBytes)); err != nil {
		return fmt.Errorf("failed to deserialize public witness: %w", err)
	}

	// Verify the proof
	return groth16.Verify(proof, vk, publicWitness)
}
