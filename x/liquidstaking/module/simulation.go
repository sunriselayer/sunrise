package liquidstaking

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunrise-zone/sunrise-app/testutil/sample"
	liquidstakingsimulation "github.com/sunrise-zone/sunrise-app/x/liquidstaking/simulation"
	"github.com/sunrise-zone/sunrise-app/x/liquidstaking/types"
)

// avoid unused import issue
var (
	_ = liquidstakingsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgLiquidStake = "op_weight_msg_liquid_stake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgLiquidStake int = 100

	opWeightMsgLiquidUnstake = "op_weight_msg_liquid_unstake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgLiquidUnstake int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	liquidstakingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&liquidstakingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgLiquidStake int
	simState.AppParams.GetOrGenerate(opWeightMsgLiquidStake, &weightMsgLiquidStake, nil,
		func(_ *rand.Rand) {
			weightMsgLiquidStake = defaultWeightMsgLiquidStake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgLiquidStake,
		liquidstakingsimulation.SimulateMsgLiquidStake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgLiquidUnstake int
	simState.AppParams.GetOrGenerate(opWeightMsgLiquidUnstake, &weightMsgLiquidUnstake, nil,
		func(_ *rand.Rand) {
			weightMsgLiquidUnstake = defaultWeightMsgLiquidUnstake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgLiquidUnstake,
		liquidstakingsimulation.SimulateMsgLiquidUnstake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgLiquidStake,
			defaultWeightMsgLiquidStake,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidstakingsimulation.SimulateMsgLiquidStake(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgLiquidUnstake,
			defaultWeightMsgLiquidUnstake,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidstakingsimulation.SimulateMsgLiquidUnstake(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
