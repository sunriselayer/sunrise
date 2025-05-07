package types

import (
	"cosmossdk.io/math"
)

const (
	// ExpansionFactor            = 2
	ShardSize                  = 4096
	ValidatorReplicationFactor = 12
)

var (
	CommitmentThresholdAmongShards = math.LegacyMustNewDecFromStr("0.666")
	CommitmentThresholdPerShard    = math.LegacyMustNewDecFromStr("0.5")
)
