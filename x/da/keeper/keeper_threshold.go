package keeper

import (
	"context"

	"cosmossdk.io/math"
)

func (k Keeper) GetZkpThreshold(ctx context.Context, shardCount uint64) uint64 {
	numActiveValidators := int64(0)
	iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		k.Logger.Error(err.Error())
		return 0
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		numActiveValidators++
	}

	// TODO: error handling
	params, _ := k.Params.Get(ctx)
	replicationFactor := math.LegacyMustNewDecFromStr(params.ReplicationFactor) // TODO: remove with Dec
	threshold := replicationFactor.MulInt64(int64(shardCount)).QuoInt64(int64(numActiveValidators)).RoundInt64()

	if threshold < 1 {
		threshold = 1
	}

	if threshold > int64(shardCount) {
		threshold = int64(shardCount)
	}

	return uint64(threshold)
}
