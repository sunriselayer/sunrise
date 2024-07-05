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

	opWeightMsgIncreaseLiquidity = "op_weight_msg_position"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePosition int = 100

	opWeightMsgDecreaseLiquidity = "op_weight_msg_position"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePosition int = 100

	opWeightMsgClaimRewards = "op_weight_msg_claim_rewards"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimRewards int = 100

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
		Pools: []types.Pool{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		PoolCount: 2,
		Positions: []types.Position{
			{
				Id:      0,
				Address: sample.AccAddress(),
			},
			{
				Id:      1,
				Address: sample.AccAddress(),
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

	var weightMsgIncreaseLiquidity int
	simState.AppParams.GetOrGenerate(opWeightMsgIncreaseLiquidity, &weightMsgIncreaseLiquidity, nil,
		func(_ *rand.Rand) {
			weightMsgIncreaseLiquidity = defaultWeightMsgUpdatePosition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgIncreaseLiquidity,
		liquiditypoolsimulation.SimulateMsgIncreaseLiquidity(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDecreaseLiquidity int
	simState.AppParams.GetOrGenerate(opWeightMsgDecreaseLiquidity, &weightMsgDecreaseLiquidity, nil,
		func(_ *rand.Rand) {
			weightMsgDecreaseLiquidity = defaultWeightMsgDeletePosition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDecreaseLiquidity,
		liquiditypoolsimulation.SimulateMsgDecreaseLiquidity(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimRewards int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimRewards, &weightMsgClaimRewards, nil,
		func(_ *rand.Rand) {
			weightMsgClaimRewards = defaultWeightMsgClaimRewards
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimRewards,
		liquiditypoolsimulation.SimulateMsgClaimRewards(am.accountKeeper, am.bankKeeper, am.keeper),
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
			opWeightMsgCreatePosition,
			defaultWeightMsgCreatePosition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgCreatePosition(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgIncreaseLiquidity,
			defaultWeightMsgUpdatePosition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgIncreaseLiquidity(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDecreaseLiquidity,
			defaultWeightMsgDeletePosition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgDecreaseLiquidity(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgClaimRewards,
			defaultWeightMsgClaimRewards,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				liquiditypoolsimulation.SimulateMsgClaimRewards(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
