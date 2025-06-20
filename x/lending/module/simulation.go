package lending

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	lendingsimulation "github.com/sunriselayer/sunrise/x/lending/simulation"
	"github.com/sunriselayer/sunrise/x/lending/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	lendingGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&lendingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgSupply          = "op_weight_msg_lending"
		defaultWeightMsgSupply int = 100
	)

	var weightMsgSupply int
	simState.AppParams.GetOrGenerate(opWeightMsgSupply, &weightMsgSupply, nil,
		func(_ *rand.Rand) {
			weightMsgSupply = defaultWeightMsgSupply
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSupply,
		lendingsimulation.SimulateMsgSupply(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgBorrow          = "op_weight_msg_lending"
		defaultWeightMsgBorrow int = 100
	)

	var weightMsgBorrow int
	simState.AppParams.GetOrGenerate(opWeightMsgBorrow, &weightMsgBorrow, nil,
		func(_ *rand.Rand) {
			weightMsgBorrow = defaultWeightMsgBorrow
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBorrow,
		lendingsimulation.SimulateMsgBorrow(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRepay          = "op_weight_msg_lending"
		defaultWeightMsgRepay int = 100
	)

	var weightMsgRepay int
	simState.AppParams.GetOrGenerate(opWeightMsgRepay, &weightMsgRepay, nil,
		func(_ *rand.Rand) {
			weightMsgRepay = defaultWeightMsgRepay
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRepay,
		lendingsimulation.SimulateMsgRepay(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgLiquidate          = "op_weight_msg_lending"
		defaultWeightMsgLiquidate int = 100
	)

	var weightMsgLiquidate int
	simState.AppParams.GetOrGenerate(opWeightMsgLiquidate, &weightMsgLiquidate, nil,
		func(_ *rand.Rand) {
			weightMsgLiquidate = defaultWeightMsgLiquidate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgLiquidate,
		lendingsimulation.SimulateMsgLiquidate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
