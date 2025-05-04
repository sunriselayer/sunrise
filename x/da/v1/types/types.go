package types

import "math"

func CalculateShardCount(blobSize uint64) uint32 {
	return ExpansionFactor * uint32(math.Ceil(float64(blobSize)/float64(ShardSize)))
}

func CalculateShardCountPerValidator(shardCount uint32, validatorCount uint32) uint32 {
	return uint32(math.Min(
		float64(shardCount),
		math.Ceil(float64(ValidatorReplicationFactor)*float64(shardCount)/float64(validatorCount)),
	))
}
