package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
)

const (
	// ModuleName defines the module name
	ModuleName = "da"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// should be changed to use collections
var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("params/")

	CommitmentKeysKeyPrefix             = collections.NewPrefix("commitment_keys/")
	BlobDeclarationsKeyPrefix           = collections.NewPrefix("blob_declarations/")
	BlobDeclarationsByExpiryPrefix      = collections.NewPrefix("blob_declarations_by_expiry/")
	BlobDeclarationsByBlockHeightPrefix = collections.NewPrefix("blob_declarations_by_block_height/")
	ValidatorsPowerSnapshotsKeyPrefix   = collections.NewPrefix("validators_power_snapshots/")
	BlobCommitmentsKeyPrefix            = collections.NewPrefix("blob_commitments/")
	BlobCommitmentsByExpiryPrefix       = collections.NewPrefix("blob_commitments_by_expiry/")
	ChallengesKeyPrefix                 = collections.NewPrefix("challenges/")
	ChallengesByShardsMerkleRootPrefix  = collections.NewPrefix("challenges_by_shards_merkle_root/")
	ChallengeIdKeyPrefix                = collections.NewPrefix("challenge_id/")
)

var (
	CommitmentKeyCodec              = collections.BytesKey
	BlobDeclarationKeyCodec         = collections.BytesKey
	ValidatorsPowerSnapshotKeyCodec = collections.Int64Key
	BlobCommitmentKeyCodec          = collections.BytesKey
	ChallengeKeyCodec               = collections.Uint64Key
)

type BlobDeclarationIndexes struct {
	Expiry      *indexes.Multi[int64, []byte, BlobDeclaration]
	BlockHeight *indexes.Multi[int64, []byte, BlobDeclaration]
}

func (i BlobDeclarationIndexes) IndexesList() []collections.Index[[]byte, BlobDeclaration] {
	return []collections.Index[[]byte, BlobDeclaration]{
		i.Expiry,
		i.BlockHeight,
	}
}

func NewBlobDeclarationIndexes(sb *collections.SchemaBuilder) BlobDeclarationIndexes {
	return BlobDeclarationIndexes{
		Expiry: indexes.NewMulti(
			sb,
			BlobDeclarationsByExpiryPrefix,
			"blob_declaration_by_expiry",
			collections.Int64Key,
			BlobDeclarationKeyCodec,
			func(_ []byte, v BlobDeclaration) (int64, error) {
				return v.Expiry.Unix(), nil
			},
		),
		BlockHeight: indexes.NewMulti(
			sb,
			BlobDeclarationsByBlockHeightPrefix,
			"blob_declaration_by_block_height",
			collections.Int64Key,
			BlobDeclarationKeyCodec,
			func(_ []byte, v BlobDeclaration) (int64, error) {
				return v.BlockHeight, nil
			},
		),
	}
}

type BlobCommitmentIndexes struct {
	Expiry *indexes.Multi[int64, []byte, BlobCommitment]
}

func (i BlobCommitmentIndexes) IndexesList() []collections.Index[[]byte, BlobCommitment] {
	return []collections.Index[[]byte, BlobCommitment]{
		i.Expiry,
	}
}

func NewBlobCommitmentIndexes(sb *collections.SchemaBuilder) BlobCommitmentIndexes {
	return BlobCommitmentIndexes{
		Expiry: indexes.NewMulti(
			sb,
			BlobCommitmentsByExpiryPrefix,
			"blob_commitment_by_expiry",
			collections.Int64Key,
			BlobCommitmentKeyCodec,
			func(_ []byte, v BlobCommitment) (int64, error) {
				return v.Expiry.Unix(), nil
			},
		),
	}
}

type ChallengeIndexes struct {
	ShardsMerkleRoot *indexes.Multi[collections.Pair[[]byte, uint64], uint64, Challenge]
}

func (i ChallengeIndexes) IndexesList() []collections.Index[uint64, Challenge] {
	return []collections.Index[uint64, Challenge]{
		i.ShardsMerkleRoot,
	}
}

func NewChallengeIndexes(sb *collections.SchemaBuilder) ChallengeIndexes {
	return ChallengeIndexes{
		ShardsMerkleRoot: indexes.NewMulti(
			sb,
			ChallengesByShardsMerkleRootPrefix,
			"challenges_by_shards_merkle_root",
			collections.PairKeyCodec(collections.BytesKey, collections.Uint64Key),
			ChallengeKeyCodec,
			func(_ uint64, v Challenge) (collections.Pair[[]byte, uint64], error) {
				return collections.Join(v.ShardsMerkleRoot, v.Id), nil
			},
		),
	}
}
