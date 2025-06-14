package ybtbrand

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	ybtbrandsimulation "github.com/sunriselayer/sunrise/x/ybtbrand/simulation"
	"github.com/sunriselayer/sunrise/x/ybtbrand/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ybtbrandGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ybtbrandGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreate          = "op_weight_msg_ybtbrand"
		defaultWeightMsgCreate int = 100
	)

	var weightMsgCreate int
	simState.AppParams.GetOrGenerate(opWeightMsgCreate, &weightMsgCreate, nil,
		func(_ *rand.Rand) {
			weightMsgCreate = defaultWeightMsgCreate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreate,
		ybtbrandsimulation.SimulateMsgCreate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgMint          = "op_weight_msg_ybtbrand"
		defaultWeightMsgMint int = 100
	)

	var weightMsgMint int
	simState.AppParams.GetOrGenerate(opWeightMsgMint, &weightMsgMint, nil,
		func(_ *rand.Rand) {
			weightMsgMint = defaultWeightMsgMint
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMint,
		ybtbrandsimulation.SimulateMsgMint(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgBurn          = "op_weight_msg_ybtbrand"
		defaultWeightMsgBurn int = 100
	)

	var weightMsgBurn int
	simState.AppParams.GetOrGenerate(opWeightMsgBurn, &weightMsgBurn, nil,
		func(_ *rand.Rand) {
			weightMsgBurn = defaultWeightMsgBurn
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurn,
		ybtbrandsimulation.SimulateMsgBurn(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAddYields          = "op_weight_msg_ybtbrand"
		defaultWeightMsgAddYields int = 100
	)

	var weightMsgAddYields int
	simState.AppParams.GetOrGenerate(opWeightMsgAddYields, &weightMsgAddYields, nil,
		func(_ *rand.Rand) {
			weightMsgAddYields = defaultWeightMsgAddYields
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddYields,
		ybtbrandsimulation.SimulateMsgAddYields(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgClaimYields          = "op_weight_msg_ybtbrand"
		defaultWeightMsgClaimYields int = 100
	)

	var weightMsgClaimYields int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimYields, &weightMsgClaimYields, nil,
		func(_ *rand.Rand) {
			weightMsgClaimYields = defaultWeightMsgClaimYields
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimYields,
		ybtbrandsimulation.SimulateMsgClaimYields(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateAdmin          = "op_weight_msg_ybtbrand"
		defaultWeightMsgUpdateAdmin int = 100
	)

	var weightMsgUpdateAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAdmin, &weightMsgUpdateAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAdmin = defaultWeightMsgUpdateAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAdmin,
		ybtbrandsimulation.SimulateMsgUpdateAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
