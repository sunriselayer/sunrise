package metadata

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"

	"github.com/sunriselayer/sunrise/x/da/das/kzg"
	"github.com/sunriselayer/sunrise/x/da/das/types"
)

func GenerateOpeningProofs(coeffs [][]fr.Element, indices []types.OpeningProofIndex) ([]types.OpeningProof, error) {
	proofs := make([]types.OpeningProof, len(indices))

	for i, index := range indices {
		proof, err := kzg.KzgOpen(coeffs[index.RowIndex], int(index.ColIndex))
		if err != nil {
			return nil, err
		}
		proofs[i] = proof
	}

	return proofs, nil
}
