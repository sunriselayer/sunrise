package liquiditypool

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
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
		k.SetAccumulatorPosition(ctx, elem.Name, elem.AccumValuePerShare, elem.Index, elem.NumShares, elem.UnclaimedRewardsTotal)
	}

	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.Pools = k.GetAllPools(ctx)
	genesis.PoolCount = k.GetPoolCount(ctx)
	genesis.Positions = k.GetAllPositions(ctx)
	genesis.PositionCount = k.GetPositionCount(ctx)
	genesis.Accumulators = k.GetAllAccumulators(ctx)
	genesis.AccumulatorPositions = k.GetAllAccumulatorPositions(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
