package lockup

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/simsx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/sunriselayer/sunrise/testutil/sample"
	"github.com/sunriselayer/sunrise/x/lockup/simulation"
	"github.com/sunriselayer/sunrise/x/lockup/types"
)

// avoid unused import issue
var (
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.AddressBech32
	}
	lockupGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&lockupGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// ProposalMsgsX returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgsX(weights simsx.WeightSource, reg simsx.Registry) {
	reg.Add(weights.Get("msg_update_params", 100), simulation.MsgUpdateParamsFactory())
}

// WeightedOperationsX returns the all the module operations with their respective weights.
func (am AppModule) WeightedOperationsX(weights simsx.WeightSource, reg simsx.Registry) {
	reg.Add(weights.Get("msg__init_lockup_account", 100), simulation.MsgInitLockupAccountFactory(am.keeper))
	reg.Add(weights.Get("msg__non_voting_delegate", 100), simulation.MsgNonVotingDelegateFactory(am.keeper))
	reg.Add(weights.Get("msg__non_voting_undelegate", 100), simulation.MsgNonVotingUndelegateFactory(am.keeper))
	reg.Add(weights.Get("msg__claim_rewards", 100), simulation.MsgClaimRewardsFactory(am.keeper))
	reg.Add(weights.Get("msg__send", 100), simulation.MsgSendFactory(am.keeper))

}
