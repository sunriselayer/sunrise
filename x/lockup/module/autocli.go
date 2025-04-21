package lockup

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/lockup/types"
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
					RpcMethod: "InitLockupAccount",
					Use:       "init-lockup-account [owner] [start_time] [end_time] [amount] [denom]",
					Short:     "Send a init-lockup-account tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "start_time"},
						{ProtoField: "end_time"},
						{ProtoField: "amount.amount"},
						{ProtoField: "amount.denom"},
					},
				},
				{
					RpcMethod: "NonVotingDelegate",
					Use:       "non-voting-delegate [id] [validator_address] [amount] [denom]",
					Short:     "Send a non-voting-delegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
						{ProtoField: "validator_address"},
						{ProtoField: "amount.amount"},
						{ProtoField: "amount.denom"},
					},
				},
				{
					RpcMethod: "NonVotingUndelegate",
					Use:       "non-voting-undelegate [id] [validator_address] [amount] [denom]",
					Short:     "Send a non-voting-undelegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
						{ProtoField: "validator_address"},
						{ProtoField: "amount.amount"},
						{ProtoField: "amount.denom"},
					},
				},
				{
					RpcMethod: "ClaimRewards",
					Use:       "claim-rewards [id]",
					Short:     "Send a claim-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "Send",
					Use:       "send [id] [recipient] [amount]",
					Short:     "Send a send tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
						{ProtoField: "recipient"},
						{ProtoField: "amount", Varargs: true},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
