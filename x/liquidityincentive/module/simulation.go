package liquidityincentive

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunriselayer/sunrise/testutil/sample"
	liquidityincentivesimulation "github.com/sunriselayer/sunrise/x/liquidityincentive/simulation"
	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// avoid unused import issue
var (
	_ = liquidityincentivesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgVoteGauge = "op_weight_msg_vote_gauge"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVoteGauge int = 100

	opWeightMsgCollectIncentiveRewards = "op_weight_msg_collect_incentive_rewards"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCollectIncentiveRewards int = 100

	opWeightMsgCollectVoteRewards = "op_weight_msg_collect_vote_rewards"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCollectVoteRewards int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	liquidityincentiveGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&liquidityincentiveGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgVoteGauge int
	simState.AppParams.GetOrGenerate(opWeightMsgVoteGauge, &weightMsgVoteGauge, nil,
		func(_ *rand.Rand) {
			weightMsgVoteGauge = defaultWeightMsgVoteGauge
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVoteGauge,
		liquidityincentivesimulation.SimulateMsgVoteGauge(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCollectIncentiveRewards int
	simState.AppParams.GetOrGenerate(opWeightMsgCollectIncentiveRewards, &weightMsgCollectIncentiveRewards, nil,
		func(_ *rand.Rand) {
			weightMsgCollectIncentiveRewards = defaultWeightMsgCollectIncentiveRewards
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCollectIncentiveRewards,
		liquidityincentivesimulation.SimulateMsgCollectIncentiveRewards(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCollectVoteRewards int
	simState.AppParams.GetOrGenerate(opWeightMsgCollectVoteRewards, &weightMsgCollectVoteRewards, nil,
		func(_ *rand.Rand) {
			weightMsgCollectVoteRewards = defaultWeightMsgCollectVoteRewards
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCollectVoteRewards,
		liquidityincentivesimulation.SimulateMsgCollectVoteRewards(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgVoteGauge,
			defaultWeightMsgVoteGauge,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidityincentivesimulation.SimulateMsgVoteGauge(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCollectIncentiveRewards,
			defaultWeightMsgCollectIncentiveRewards,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidityincentivesimulation.SimulateMsgCollectIncentiveRewards(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCollectVoteRewards,
			defaultWeightMsgCollectVoteRewards,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquidityincentivesimulation.SimulateMsgCollectVoteRewards(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
