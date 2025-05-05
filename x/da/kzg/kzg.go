package kzg

import (
	"fmt"
	"runtime"

	fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	kzg "github.com/consensys/gnark-crypto/ecc/bls12-381/kzg"
)

func xs() []fr.Element {
	xs := make([]fr.Element, EvaluationPointCount)
	for i := range xs {
		xs[i].SetUint64(uint64(i))
	}
	return xs
}

// ys converts 1024 bytes of data into 32 fr.Element points
// Each point is 32 bytes (256 bits)
func ys(data []byte) ([]fr.Element, error) {
	if len(data) != DataSize {
		return nil, fmt.Errorf("data must be exactly %d bytes", DataSize)
	}

	points := make([]fr.Element, EvaluationPointCount)
	for i := range points {
		// Convert each 32-byte chunk to fr.Element
		points[i].SetBytes(data[i*32 : (i+1)*32])
	}
	return points, nil
}

// KzgCommit generates a KZG commitment for the given data
// data: original 1024 bytes of data
// srs: SRS containing proving key
func KzgCommit(data []byte, srs kzg.SRS) ([]byte, error) {
	xs := xs()
	ys, err := ys(data)
	if err != nil {
		return nil, err
	}

	poly := InterpolateLagrange(xs, ys)

	// Use number of CPU cores for parallel processing
	nbTasks := runtime.NumCPU()
	commitment, err := kzg.Commit(poly, srs.Pk, nbTasks)
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
func KzgOpen(pointIndex int, data []byte) (kzg.OpeningProof, fr.Element, error) {
	if pointIndex < 0 || pointIndex >= EvaluationPointCount {
		return kzg.OpeningProof{}, fr.Element{}, fmt.Errorf("pointIndex must be between 0 and 31")
	}

	xs := xs()
	ys, err := ys(data)
	if err != nil {
		return kzg.OpeningProof{}, fr.Element{}, err
	}

	poly := InterpolateLagrange(xs, ys)
	point := xs[pointIndex]
	value := ys[pointIndex]

	// Generate opening proof
	proof, err := kzg.Open(poly, point, Srs.Pk)
	if err != nil {
		return kzg.OpeningProof{}, fr.Element{}, err
	}

	return proof, value, nil
}

// KzgVerify verifies an opening proof
// commitment: the KZG commitment
// pointIndex: index of the point being proved (0-31)
// proof: the opening proof
// srs: SRS containing verification key
func KzgVerify(commitment []byte, proof kzg.OpeningProof, pointIndex int) error {
	if pointIndex < 0 || pointIndex >= EvaluationPointCount {
		return fmt.Errorf("pointIndex must be between 0 and 31")
	}

	commitmentDigest := kzg.Digest{}
	_, err := commitmentDigest.SetBytes(commitment)
	if err != nil {
		return err
	}

	// Get point from index
	point := xs()[pointIndex]

	// Verify the proof
	return kzg.Verify(&commitmentDigest, &proof, point, Srs.Vk)
}
