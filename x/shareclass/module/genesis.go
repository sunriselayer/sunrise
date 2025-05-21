package shareclass

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sunriselayer/sunrise/x/shareclass/keeper"
	"github.com/sunriselayer/sunrise/x/shareclass/types"
)

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	err := k.InitGenesis(ctx, genState)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genState, err := k.ExportGenesis(ctx)
	if err != nil {
		panic(err)
	}

	return genState
}
