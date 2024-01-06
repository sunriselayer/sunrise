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
	opWeightMsgMintDerivative = "op_weight_msg_mint_derivative"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMintDerivative int = 100

	opWeightMsgBurnDerivative = "op_weight_msg_burn_derivative"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBurnDerivative int = 100

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

	var weightMsgMintDerivative int
	simState.AppParams.GetOrGenerate(opWeightMsgMintDerivative, &weightMsgMintDerivative, nil,
		func(_ *rand.Rand) {
			weightMsgMintDerivative = defaultWeightMsgMintDerivative
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintDerivative,
		liquidstakingsimulation.SimulateMsgMintDerivative(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgBurnDerivative int
	simState.AppParams.GetOrGenerate(opWeightMsgBurnDerivative, &weightMsgBurnDerivative, nil,
		func(_ *rand.Rand) {
			weightMsgBurnDerivative = defaultWeightMsgBurnDerivative
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBurnDerivative,
		liquidstakingsimulation.SimulateMsgBurnDerivative(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgMintDerivative,
			defaultWeightMsgMintDerivative,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidstakingsimulation.SimulateMsgMintDerivative(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgBurnDerivative,
			defaultWeightMsgBurnDerivative,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidstakingsimulation.SimulateMsgBurnDerivative(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
