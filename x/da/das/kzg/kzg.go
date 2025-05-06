package kzg

import (
	"fmt"
	"runtime"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/kzg"

	"github.com/sunriselayer/sunrise/x/da/das/consts"
	"github.com/sunriselayer/sunrise/x/da/das/erasurecoding"
	"github.com/sunriselayer/sunrise/x/da/das/types"
)

func KzgCommit(coeffs []fr.Element, srs kzg.SRS) ([]byte, error) {
	// Use number of CPU cores for parallel processing
	nbTasks := runtime.NumCPU()
	commitment, err := kzg.Commit(coeffs, srs.Pk, nbTasks)
	if err != nil {
		return nil, err
	}

	// Convert fixed-size array to slice
	bytes := commitment.Bytes()
	return bytes[:], nil
}

// KzgOpen generates an opening proof for a specific point in the polynomial
// pointIndex: index of the point to prove (0-31)
// data: original 1024 bytes of data
// srs: SRS containing proving key
func KzgOpen(coeffs []fr.Element, pointIndex int) (types.OpeningProof, error) {
	if pointIndex < 0 || pointIndex >= len(coeffs) {
		return types.OpeningProof{}, fmt.Errorf("pointIndex must be between 0 and %d", len(coeffs)-1)
	}

	point := erasurecoding.EvaluationPoints(uint64(pointIndex + 1))[pointIndex]

	// Generate opening proof
	proof, err := kzg.Open(coeffs, point, consts.Srs.Pk)
	if err != nil {
		return types.OpeningProof{}, err
	}

	h := proof.H.Bytes()
	claimedValue := proof.ClaimedValue.Bytes()
	proofSerializable := types.OpeningProof{
		H:            h[:],
		ClaimedValue: claimedValue[:],
	}

	return proofSerializable, nil
}

// KzgVerify verifies an opening proof
// commitment: the KZG commitment
// pointIndex: index of the point being proved (0-31)
// proof: the opening proof
// srs: SRS containing verification key
func KzgVerify(commitment []byte, proof types.OpeningProof, pointIndex int) error {
	commitmentDigest := kzg.Digest{}
	_, err := commitmentDigest.SetBytes(commitment)
	if err != nil {
		return err
	}

	h := bls12381.G1Affine{}
	_, err = h.SetBytes(proof.H)
	if err != nil {
		return err
	}

	claimedValue := fr.Element{}
	claimedValue.SetBytes(proof.ClaimedValue)

	proofRaw := kzg.OpeningProof{
		H:            h,
		ClaimedValue: claimedValue,
	}

	// Get point from index
	point := erasurecoding.EvaluationPoints(uint64(pointIndex + 1))[pointIndex]

	// Verify the proof
	return kzg.Verify(&commitmentDigest, &proofRaw, point, consts.Srs.Vk)
}
