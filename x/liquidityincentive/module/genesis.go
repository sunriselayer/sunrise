package liquidityincentive

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/keeper"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the epoch
	for _, elem := range genState.EpochList {
		k.SetEpoch(ctx, elem)
	}

	// Set epoch count
	k.SetEpochCount(ctx, genState.EpochCount)
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.EpochList = k.GetAllEpoch(ctx)
	genesis.EpochCount = k.GetEpochCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
