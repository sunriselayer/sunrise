package liquidstaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"sunrise/x/liquidstaking/keeper"
	"sunrise/x/liquidstaking/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := types.ValidateGenesis(genState); err != nil {
		panic(err)
	}
	// init to prevent nil slice, []types.WhitelistedValidator(nil)
	if genState.Params.WhitelistedValidators == nil || len(genState.Params.WhitelistedValidators) == 0 {
		genState.Params.WhitelistedValidators = []types.WhitelistedValidator{}
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	for _, lv := range genState.LiquidValidators {
		k.SetLiquidValidator(ctx, lv)
	}

	moduleAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	// init to prevent nil slice, []types.WhitelistedValidator(nil)
	if params.WhitelistedValidators == nil || len(params.WhitelistedValidators) == 0 {
		params.WhitelistedValidators = []types.WhitelistedValidator{}
	}

	genesis.LiquidValidators := k.GetAllLiquidValidators(ctx)

	return genesis
}
