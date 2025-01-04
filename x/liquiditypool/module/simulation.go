package liquiditypool

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"sunrise/testutil/sample"
	"sunrise/x/liquiditypool/simulation"
	"sunrise/x/liquiditypool/types"
)

// avoid unused import issue
var (
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.AddressBech32
	}
	liquiditypoolGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&liquiditypoolGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// ProposalMsgsX returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgsX(weights simsx.WeightSource, reg simsx.Registry) {
	reg.Add(weights.Get("msg_update_params", 100), simulation.MsgUpdateParamsFactory())
}

// WeightedOperationsX returns the all the module operations with their respective weights.
func (am AppModule) WeightedOperationsX(weights simsx.WeightSource, reg simsx.Registry) {
	reg.Add(weights.Get("msg__create_pool", 100), simulation.MsgCreatePoolFactory(am.keeper))
	reg.Add(weights.Get("msg__create_position", 100), simulation.MsgCreatePositionFactory(am.keeper))
	reg.Add(weights.Get("msg__increase_liquidity", 100), simulation.MsgIncreaseLiquidityFactory(am.keeper))
	reg.Add(weights.Get("msg__decrease_liquidity", 100), simulation.MsgDecreaseLiquidityFactory(am.keeper))
	reg.Add(weights.Get("msg__claim_rewards", 100), simulation.MsgClaimRewardsFactory(am.keeper))

}
