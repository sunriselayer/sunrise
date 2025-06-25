package keeper

import (
	"context"

	"cosmossdk.io/math"
)

func (k Keeper) GetZkpThreshold(ctx context.Context, shardCount uint64) (uint64, error) {
	numActiveValidators := int64(0)
	iterator, err := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return 0, err
	}

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.StakingKeeper.Validator(ctx, iterator.Value())
		if err != nil {
			k.Logger().Error(err.Error())
			continue
		}
		if validator.IsBonded() {
			numActiveValidators++
		}
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return 0, err
	}
	replicationFactor := math.LegacyMustNewDecFromStr(params.ReplicationFactor) // TODO: remove with Dec
	if numActiveValidators == 0 {
		return 0, nil
	}
	threshold := min(max(replicationFactor.MulInt64(int64(shardCount)).QuoInt64(int64(numActiveValidators)).Ceil().TruncateInt64(), 1), int64(shardCount))

	return uint64(threshold), nil
}
