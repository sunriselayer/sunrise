package liquiditypool

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sunriselayer/sunrise/testutil/sample"
	liquiditypoolsimulation "github.com/sunriselayer/sunrise/x/liquiditypool/simulation"
	"github.com/sunriselayer/sunrise/x/liquiditypool/types"
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

	opWeightMsgCreatePosition = "op_weight_msg_position"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePosition int = 100

	opWeightMsgUpdatePosition = "op_weight_msg_position"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePosition int = 100

	opWeightMsgDeletePosition = "op_weight_msg_position"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePosition int = 100

	opWeightMsgCollectFees = "op_weight_msg_collect_fees"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCollectFees int = 100

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
				Id:     0,
				Sender: sample.AccAddress(),
			},
			{
				Id:     1,
				Sender: sample.AccAddress(),
			},
		},
		PoolCount: 2,
		PositionList: []types.Position{
			{
				Id:     0,
				Sender: sample.AccAddress(),
			},
			{
				Id:     1,
				Sender: sample.AccAddress(),
			},
		},
		PositionCount: 2,
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

	var weightMsgCreatePosition int
	simState.AppParams.GetOrGenerate(opWeightMsgCreatePosition, &weightMsgCreatePosition, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePosition = defaultWeightMsgCreatePosition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePosition,
		liquiditypoolsimulation.SimulateMsgCreatePosition(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePosition int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdatePosition, &weightMsgUpdatePosition, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePosition = defaultWeightMsgUpdatePosition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePosition,
		liquiditypoolsimulation.SimulateMsgUpdatePosition(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePosition int
	simState.AppParams.GetOrGenerate(opWeightMsgDeletePosition, &weightMsgDeletePosition, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePosition = defaultWeightMsgDeletePosition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeletePosition,
		liquiditypoolsimulation.SimulateMsgDeletePosition(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCollectFees int
	simState.AppParams.GetOrGenerate(opWeightMsgCollectFees, &weightMsgCollectFees, nil,
		func(_ *rand.Rand) {
			weightMsgCollectFees = defaultWeightMsgCollectFees
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCollectFees,
		liquiditypoolsimulation.SimulateMsgCollectFees(am.accountKeeper, am.bankKeeper, am.keeper),
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
			opWeightMsgCreatePosition,
			defaultWeightMsgCreatePosition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgCreatePosition(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdatePosition,
			defaultWeightMsgUpdatePosition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgUpdatePosition(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeletePosition,
			defaultWeightMsgDeletePosition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgDeletePosition(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCollectFees,
			defaultWeightMsgCollectFees,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgCollectFees(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
