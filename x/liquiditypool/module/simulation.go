package liquiditypool

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunriselayer/sunrise-app/testutil/sample"
	liquiditypoolsimulation "github.com/sunriselayer/sunrise-app/x/liquiditypool/simulation"
	"github.com/sunriselayer/sunrise-app/x/liquiditypool/types"
)

// avoid unused import issue
var (
	_ = liquiditypoolsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
opWeightMsgCreatePool = "op_weight_msg_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePool int = 100

	opWeightMsgUpdatePool = "op_weight_msg_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePool int = 100

	opWeightMsgDeletePool = "op_weight_msg_pool"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePool int = 100

	opWeightMsgCreateTwap = "op_weight_msg_twap"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateTwap int = 100

	opWeightMsgUpdateTwap = "op_weight_msg_twap"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateTwap int = 100

	opWeightMsgDeleteTwap = "op_weight_msg_twap"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteTwap int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	liquiditypoolGenesis := types.GenesisState{
		Params: types.DefaultParams(),
			PoolList: []types.Pool{
		{
			Id: 0,
			Creator: sample.AccAddress(),

		},
		{
			Id: 1,
			Creator: sample.AccAddress(),

		},
	},
	PoolCount: 2,
	TwapList: []types.Twap{
		{
			Creator: sample.AccAddress(),
Index: "0",
},
		{
			Creator: sample.AccAddress(),
Index: "1",
},
	},
	// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&liquiditypoolGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreatePool int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = defaultWeightMsgCreatePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePool,
		liquiditypoolsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePool int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePool, &weightMsgUpdatePool, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePool = defaultWeightMsgUpdatePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePool,
		liquiditypoolsimulation.SimulateMsgUpdatePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePool int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePool, &weightMsgDeletePool, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePool = defaultWeightMsgDeletePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePool,
		liquiditypoolsimulation.SimulateMsgDeletePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateTwap int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTwap, &weightMsgCreateTwap, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTwap = defaultWeightMsgCreateTwap
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTwap,
		liquiditypoolsimulation.SimulateMsgCreateTwap(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateTwap int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTwap, &weightMsgUpdateTwap, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTwap = defaultWeightMsgUpdateTwap
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTwap,
		liquiditypoolsimulation.SimulateMsgUpdateTwap(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteTwap int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteTwap, &weightMsgDeleteTwap, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteTwap = defaultWeightMsgDeleteTwap
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteTwap,
		liquiditypoolsimulation.SimulateMsgDeleteTwap(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
	opWeightMsgCreatePool,
	defaultWeightMsgCreatePool,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		liquiditypoolsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgUpdatePool,
	defaultWeightMsgUpdatePool,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		liquiditypoolsimulation.SimulateMsgUpdatePool(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgDeletePool,
	defaultWeightMsgDeletePool,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		liquiditypoolsimulation.SimulateMsgDeletePool(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgCreateTwap,
	defaultWeightMsgCreateTwap,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		liquiditypoolsimulation.SimulateMsgCreateTwap(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgUpdateTwap,
	defaultWeightMsgUpdateTwap,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		liquiditypoolsimulation.SimulateMsgUpdateTwap(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgDeleteTwap,
	defaultWeightMsgDeleteTwap,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		liquiditypoolsimulation.SimulateMsgDeleteTwap(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
