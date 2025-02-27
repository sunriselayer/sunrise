package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set all the pool
	k.SetPoolCount(ctx, genState.PoolCount)
	for _, elem := range genState.Pools {
		k.SetPool(ctx, elem)
	}

	// Set all the position
	k.SetPositionCount(ctx, genState.PositionCount)
	for _, elem := range genState.Positions {
		k.SetPosition(ctx, elem)
	}

	// Set all accumulators
	for _, elem := range genState.Accumulators {
		err := k.SetAccumulator(ctx, elem)
		if err != nil {
			panic(err)
		}
	}
	// Set all accumulator positions
	for _, elem := range genState.AccumulatorPositions {
		numShares, err := math.LegacyNewDecFromStr(elem.NumShares)
		if err != nil {
			panic(err)
		}
		k.SetAccumulatorPosition(ctx, elem.Name, elem.AccumValuePerShare, elem.Index, numShares, elem.UnclaimedRewardsTotal)
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	genesis.Pools = k.GetAllPools(ctx)
	genesis.PoolCount = k.GetPoolCount(ctx)
	genesis.Positions = k.GetAllPositions(ctx)
	genesis.PositionCount = k.GetPositionCount(ctx)
	genesis.Accumulators = k.GetAllAccumulators(ctx)
	genesis.AccumulatorPositions = k.GetAllAccumulatorPositions(ctx)

	return genesis, nil
}
