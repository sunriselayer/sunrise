package tokenconverter

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunriselayer/sunrise/testutil/sample"
	tokenconvertersimulation "github.com/sunriselayer/sunrise/x/tokenconverter/simulation"
	"github.com/sunriselayer/sunrise/x/tokenconverter/types"
)

// avoid unused import issue
var (
	_ = tokenconvertersimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgConvertExactAmountIn = "op_weight_msg_convert_exact_amount_in"
	// TODO: Determine the simulation weight value
	defaultWeightMsgConvertExactAmountIn int = 100

	opWeightMsgConvertExactAmountOut = "op_weight_msg_convert_exact_amount_out"
	// TODO: Determine the simulation weight value
	defaultWeightMsgConvertExactAmountOut int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenconverterGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenconverterGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgConvertExactAmountIn int
	simState.AppParams.GetOrGenerate(opWeightMsgConvertExactAmountIn, &weightMsgConvertExactAmountIn, nil,
		func(_ *rand.Rand) {
			weightMsgConvertExactAmountIn = defaultWeightMsgConvertExactAmountIn
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgConvertExactAmountIn,
		tokenconvertersimulation.SimulateMsgConvertExactAmountIn(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgConvertExactAmountOut int
	simState.AppParams.GetOrGenerate(opWeightMsgConvertExactAmountOut, &weightMsgConvertExactAmountOut, nil,
		func(_ *rand.Rand) {
			weightMsgConvertExactAmountOut = defaultWeightMsgConvertExactAmountOut
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgConvertExactAmountOut,
		tokenconvertersimulation.SimulateMsgConvertExactAmountOut(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgConvertExactAmountIn,
			defaultWeightMsgConvertExactAmountIn,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tokenconvertersimulation.SimulateMsgConvertExactAmountIn(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgConvertExactAmountOut,
			defaultWeightMsgConvertExactAmountOut,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tokenconvertersimulation.SimulateMsgConvertExactAmountOut(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
