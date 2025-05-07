package metadata

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"

	"github.com/sunriselayer/sunrise/x/da/das/consts"
	"github.com/sunriselayer/sunrise/x/da/das/erasurecoding"
	"github.com/sunriselayer/sunrise/x/da/das/kzg"
	"github.com/sunriselayer/sunrise/x/da/das/merkle"
	"github.com/sunriselayer/sunrise/x/da/das/types"
)

func GenerateMetadata(blob []byte, upload func(shard []byte) (uri string, err error)) (types.Metadata, error) {
	elements := erasurecoding.ConvertToElements(blob)
	rowElements, rowsLen, colsLen := erasurecoding.SplitElementsIntoRows(elements)

	rows := make([]types.Row, rowsLen)
	coeffs := make([][]fr.Element, rowsLen)

	for i, row := range rowElements {
		coeffs[i] = erasurecoding.CalculateCoefficients(row)
		kzgCommitment, err := kzg.KzgCommit(coeffs[i], consts.Srs)
		if err != nil {
			return types.Metadata{}, err
		}
		rows[i].KzgCommitment = kzgCommitment

		points, err := erasurecoding.CalculateExtendedPoints(coeffs[i])
		if err != nil {
			return types.Metadata{}, err
		}
		rowElements[i] = points

		shardsLen := consts.CalculateShardCountPerRow(consts.ExtensionRatio * colsLen)
		rows[i].Shards = make([]types.Shard, shardsLen)
		shardHashes := make([][32]byte, shardsLen)

		for j := range rows[i].Shards {
			shardBytes := []byte{} // TODO

			shardHashes[j] = merkle.Hash(shardBytes)
			rows[i].Shards[j].Hash = shardHashes[j][:]

			rows[i].Shards[j].Uri, err = upload(shardBytes)
			if err != nil {
				return types.Metadata{}, err
			}
		}
		root := merkle.MerkleRoot(shardHashes)
		rows[i].ShardsMerkleRoot = root[:]
	}

	indices, err := CalculateOpeningProofIndices(rowsLen, consts.ExtensionRatio*colsLen)
	if err != nil {
		return types.Metadata{}, err
	}
	proofs, err := GenerateOpeningProofs(coeffs, indices)
	if err != nil {
		return types.Metadata{}, err
	}

	metadata := types.Metadata{
		Rows:                rows,
		OpeningProofIndices: indices,
		OpeningProofs:       proofs,
	}

	return metadata, nil
}
