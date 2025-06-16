package ybtbase

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	ybtbasesimulation "github.com/sunriselayer/sunrise/x/ybtbase/simulation"
	"github.com/sunriselayer/sunrise/x/ybtbase/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	ybtbaseGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&ybtbaseGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreate          = "op_weight_msg_ybtbase"
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
		ybtbasesimulation.SimulateMsgCreate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgMint          = "op_weight_msg_ybtbase"
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
		ybtbasesimulation.SimulateMsgMint(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgBurn          = "op_weight_msg_ybtbase"
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
		ybtbasesimulation.SimulateMsgBurn(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgAddYield          = "op_weight_msg_ybtbase"
		defaultWeightMsgAddYield int = 100
	)

	var weightMsgAddYield int
	simState.AppParams.GetOrGenerate(opWeightMsgAddYield, &weightMsgAddYield, nil,
		func(_ *rand.Rand) {
			weightMsgAddYield = defaultWeightMsgAddYield
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddYield,
		ybtbasesimulation.SimulateMsgAddYield(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgGrantPermission          = "op_weight_msg_grant_permission"
		defaultWeightMsgGrantPermission int = 100
	)

	var weightMsgGrantPermission int
	simState.AppParams.GetOrGenerate(opWeightMsgGrantPermission, &weightMsgGrantPermission, nil,
		func(_ *rand.Rand) {
			weightMsgGrantPermission = defaultWeightMsgGrantPermission
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgGrantPermission,
		ybtbasesimulation.SimulateMsgGrantPermission(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRevokePermission          = "op_weight_msg_revoke_permission"
		defaultWeightMsgRevokePermission int = 100
	)

	var weightMsgRevokePermission int
	simState.AppParams.GetOrGenerate(opWeightMsgRevokePermission, &weightMsgRevokePermission, nil,
		func(_ *rand.Rand) {
			weightMsgRevokePermission = defaultWeightMsgRevokePermission
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRevokePermission,
		ybtbasesimulation.SimulateMsgRevokePermission(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgClaimYield          = "op_weight_msg_ybtbase"
		defaultWeightMsgClaimYield int = 100
	)

	var weightMsgClaimYield int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimYield, &weightMsgClaimYield, nil,
		func(_ *rand.Rand) {
			weightMsgClaimYield = defaultWeightMsgClaimYield
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimYield,
		ybtbasesimulation.SimulateMsgClaimYield(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateAdmin          = "op_weight_msg_ybtbase"
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
		ybtbasesimulation.SimulateMsgUpdateAdmin(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgSend          = "op_weight_msg_send"
		defaultWeightMsgSend int = 100
	)

	var weightMsgSend int
	simState.AppParams.GetOrGenerate(opWeightMsgSend, &weightMsgSend, nil,
		func(_ *rand.Rand) {
			weightMsgSend = defaultWeightMsgSend
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSend,
		ybtbasesimulation.SimulateMsgSend(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
