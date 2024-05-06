package grant

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunriselayer/sunrise/testutil/sample"
	grantsimulation "github.com/sunriselayer/sunrise/x/blobgrant/simulation"
	"github.com/sunriselayer/sunrise/x/blobgrant/types"
)

// avoid unused import issue
var (
	_ = grantsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateRegistration = "op_weight_msg_registration"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateRegistration int = 100

	opWeightMsgUpdateRegistration = "op_weight_msg_registration"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateRegistration int = 100

	opWeightMsgDeleteRegistration = "op_weight_msg_registration"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteRegistration int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	grantGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		RegistrationList: []types.Registration{
			{
				LiquidityProvider: sample.AccAddress(),
				Grantee:           "0",
			},
			{
				LiquidityProvider: sample.AccAddress(),
				Grantee:           "1",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&grantGenesis)
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

	var weightMsgCreateRegistration int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateRegistration, &weightMsgCreateRegistration, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRegistration = defaultWeightMsgCreateRegistration
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateRegistration,
		grantsimulation.SimulateMsgCreateRegistration(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateRegistration int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateRegistration, &weightMsgUpdateRegistration, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateRegistration = defaultWeightMsgUpdateRegistration
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateRegistration,
		grantsimulation.SimulateMsgUpdateRegistration(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteRegistration int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteRegistration, &weightMsgDeleteRegistration, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteRegistration = defaultWeightMsgDeleteRegistration
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteRegistration,
		grantsimulation.SimulateMsgDeleteRegistration(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateRegistration,
			defaultWeightMsgCreateRegistration,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				grantsimulation.SimulateMsgCreateRegistration(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateRegistration,
			defaultWeightMsgUpdateRegistration,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				grantsimulation.SimulateMsgUpdateRegistration(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteRegistration,
			defaultWeightMsgDeleteRegistration,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				grantsimulation.SimulateMsgDeleteRegistration(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
