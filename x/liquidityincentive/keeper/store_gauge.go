package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// SetGauge set a specific gauge in the store from its index
func (k Keeper) SetGauge(ctx context.Context, gauge types.Gauge) error {
	err := k.Gauges.Set(ctx, types.GaugeKey(gauge.PreviousEpochId, gauge.PoolId), gauge)
	if err != nil {
		return err
	}
	return nil
}

// GetGauge returns a gauge from its index
func (k Keeper) GetGauge(ctx context.Context, previousEpochId uint64, poolId uint64) (val types.Gauge, found bool, err error) {
	key := types.GaugeKey(previousEpochId, poolId)
	has, err := k.Gauges.Has(ctx, key)
	if err != nil {
		return val, false, err
	}

	if !has {
		return val, false, nil
	}

	val, err = k.Gauges.Get(ctx, key)
	if err != nil {
		return val, false, err
	}

	return val, true, nil
}

// RemoveGauge removes a gauge from the store
func (k Keeper) RemoveGauge(ctx context.Context, previousEpochId uint64, poolId uint64) error {
	err := k.Gauges.Remove(ctx, types.GaugeKey(previousEpochId, poolId))
	if err != nil {
		return err
	}
	return nil
}

// GetAllGauges returns all gauges
func (k Keeper) GetAllGauges(ctx context.Context) (list []types.Gauge, err error) {
	err = k.Gauges.Walk(
		ctx,
		nil,
		func(key collections.Pair[uint64, uint64], value types.Gauge) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAllGaugeByPreviousEpochId returns all gauges by previous epoch id
func (k Keeper) GetAllGaugeByPreviousEpochId(ctx context.Context, previousEpochId uint64) (list []types.Gauge, err error) {
	err = k.Gauges.Walk(
		ctx,
		collections.NewPrefixedPairRange[uint64, uint64](previousEpochId),
		func(key collections.Pair[uint64, uint64], value types.Gauge) (bool, error) {
			list = append(list, value)

			return false, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return list, nil
}
