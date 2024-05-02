package liquiditypool

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise-app/x/liquiditypool/keeper"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the pair
for _, elem := range genState.PairList {
	k.SetPair(ctx, elem)
}
// Set all the pool
for _, elem := range genState.PoolList {
	k.SetPool(ctx, elem)
}

// Set pool count
k.SetPoolCount(ctx, genState.PoolCount)
// Set all the twap
for _, elem := range genState.TwapList {
	k.SetTwap(ctx, elem)
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

	genesis.PairList = k.GetAllPair(ctx)
genesis.PoolList = k.GetAllPool(ctx)
genesis.PoolCount = k.GetPoolCount(ctx)
genesis.TwapList = k.GetAllTwap(ctx)
// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
