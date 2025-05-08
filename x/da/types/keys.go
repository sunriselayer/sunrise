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
	ParamsKey = collections.NewPrefix("v0_params/")

	PublishedDataKeyPrefix             = collections.NewPrefix("v0_published_data/")
	PublishedDataStatusTimeIndexPrefix = collections.NewPrefix("v0_published_data_by_status_time/")
	ChallengeCountsKeyPrefix           = collections.NewPrefix("v0_challenge_counts/")
	FaultCountsKeyPrefix               = collections.NewPrefix("v0_fault_counts/")
	ProofKeyPrefix                     = collections.NewPrefix("v0_proofs/")
	InvalidityKeyPrefix                = collections.NewPrefix("v0_invalidities/")
	ProofDeputiesKeyPrefix             = collections.NewPrefix("v0_proof_deputies/")
)

var (
	PublishedDataKeyCodec = collections.StringKey
	FaultCounterKeyCodec  = collections.BytesKey
	ProofKeyCodec         = collections.PairKeyCodec(collections.StringKey, collections.BytesKey)
	InvalidityKeyCodec    = collections.PairKeyCodec(collections.StringKey, collections.BytesKey)
	ProofDeputyKeyCodec   = collections.BytesKey
)

type PublishedDataIndexes struct {
	StatusTime *indexes.Multi[collections.Pair[string, int64], string, PublishedData]
}

func (i PublishedDataIndexes) IndexesList() []collections.Index[string, PublishedData] {
	return []collections.Index[string, PublishedData]{
		i.StatusTime,
	}
}

func NewPublishedDataIndexes(sb *collections.SchemaBuilder) PublishedDataIndexes {
	return PublishedDataIndexes{
		StatusTime: indexes.NewMulti(
			sb,
			PublishedDataStatusTimeIndexPrefix,
			"published_data_by_status_time",
			collections.PairKeyCodec(collections.StringKey, collections.Int64Key),
			collections.StringKey,
			func(_ string, v PublishedData) (collections.Pair[string, int64], error) {
				return collections.Join(v.Status.String(), v.Timestamp.Unix()), nil
			},
		),
	}
}
