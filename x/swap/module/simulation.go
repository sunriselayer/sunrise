package swap

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunriselayer/sunrise/testutil/sample"
	swapsimulation "github.com/sunriselayer/sunrise/x/swap/simulation"
	"github.com/sunriselayer/sunrise/x/swap/types"
)

// avoid unused import issue
var (
	_ = swapsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgSwapExactAmountIn = "op_weight_msg_swap_exact_amount_in"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwapExactAmountIn int = 100

	opWeightMsgSwapExactAmountOut = "op_weight_msg_swap_exact_amount_out"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwapExactAmountOut int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	swapGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&swapGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgSwapExactAmountIn int
	simState.AppParams.GetOrGenerate(opWeightMsgSwapExactAmountIn, &weightMsgSwapExactAmountIn, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactAmountIn = defaultWeightMsgSwapExactAmountIn
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactAmountIn,
		swapsimulation.SimulateMsgSwapExactAmountIn(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwapExactAmountOut int
	simState.AppParams.GetOrGenerate(opWeightMsgSwapExactAmountOut, &weightMsgSwapExactAmountOut, nil,
		func(_ *rand.Rand) {
			weightMsgSwapExactAmountOut = defaultWeightMsgSwapExactAmountOut
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwapExactAmountOut,
		swapsimulation.SimulateMsgSwapExactAmountOut(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgSwapExactAmountIn,
			defaultWeightMsgSwapExactAmountIn,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				swapsimulation.SimulateMsgSwapExactAmountIn(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSwapExactAmountOut,
			defaultWeightMsgSwapExactAmountOut,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				swapsimulation.SimulateMsgSwapExactAmountOut(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
