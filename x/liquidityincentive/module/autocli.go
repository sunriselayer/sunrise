package liquidityincentive

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/liquidityincentive/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "Epochs",
					Use:       "all-epochs",
					Short:     "List all epoch",
				},
				{
					RpcMethod:      "Epoch",
					Use:            "epoch [id]",
					Short:          "Shows a epoch by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "Votes",
					Use:       "all-votes",
					Short:     "List all gauge votes",
				},
				{
					RpcMethod:      "Vote",
					Use:            "vote [address]",
					Short:          "Shows a gauge vote",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "VoteGauge",
					Use:            "vote-gauge",
					Short:          "Send a vote-gauge tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod: "RegisterBribe",
					Use:       "register-bribe [epoch_id] [pool_id] [amount]",
					Short:     "Send a register-bribe tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "epoch_id"},
						{ProtoField: "pool_id"},
						{ProtoField: "amount", Varargs: true},
					},
				},
				{
					RpcMethod: "ClaimBribes",
					Use:       "claim-bribes [bribe_id]",
					Short:     "Send a claim-bribes tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "bribe_ids"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
