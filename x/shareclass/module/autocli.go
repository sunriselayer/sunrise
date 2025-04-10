package shareclass

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/shareclass/types"
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
					RpcMethod: "NonVotingDelegate",
					Use:       "non-voting-delegate <validator_address> <amount> <denom>",
					Short:     "Send a non-voting-delegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount.amount"},
						{ProtoField: "amount.denom"},
					},
				},
				{
					RpcMethod: "NonVotingUndelegate",
					Use:       "non-voting-undelegate <validator_address> <amount> <denom> <recipient>",
					Short:     "Send a non-voting-undelegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount.amount"},
						{ProtoField: "amount.denom"},
						{ProtoField: "recipient"},
					},
				},
				{
					RpcMethod: "ClaimRewards",
					Use:       "claim-rewards <validator_address>",
					Short:     "Send a claim-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
					},
				},
				{
					RpcMethod: "CreateValidator",
					Use:       "create-validator <validator_address> <min_self_delegation> <amount> <denom> <fee_amount> <fee_denom>",
					Short:     "Send a create-validator tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "min_self_delegation"},
						{ProtoField: "amount.amount"},
						{ProtoField: "amount.denom"},
						{ProtoField: "fee.amount"},
						{ProtoField: "fee.denom"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
