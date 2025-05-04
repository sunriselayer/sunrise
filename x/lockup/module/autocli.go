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
				{
					RpcMethod: "LockupAccounts",
					Use:       "lockup-accounts [owner]",
					Short:     "Shows the lockup accounts of the owner",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
					},
				},
				{
					RpcMethod: "LockupAccount",
					Use:       "lockup-account [owner] [lockup_account_id]",
					Short:     "Shows the lockup account of the owner and id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "lockup_account_id"},
					},
				},
				{
					RpcMethod: "SpendableAmount",
					Use:       "spendable-amount [owner] [lockup_account_id]",
					Short:     "Shows the spendable amount of the lockup account",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "lockup_account_id"},
					},
				},
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
					Use:       "init-lockup-account [owner] [start_time] [end_time] [amount]",
					Short:     "Send a init-lockup-account tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "start_time"},
						{ProtoField: "end_time"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "NonVotingDelegate",
					Use:       "non-voting-delegate [lockup_account_id] [validator_address] [amount]",
					Short:     "Send a non-voting-delegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "lockup_account_id"},
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "NonVotingUndelegate",
					Use:       "non-voting-undelegate [lockup_account_id] [validator_address] [amount]",
					Short:     "Send a non-voting-undelegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "lockup_account_id"},
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "ClaimRewards",
					Use:       "claim-rewards [lockup_account_id] [validator_address]",
					Short:     "Send a claim-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "lockup_account_id"},
						{ProtoField: "validator_address"},
					},
				},
				{
					RpcMethod: "Send",
					Use:       "send [lockup_account_id] [recipient] [amount]",
					Short:     "Send a send tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "lockup_account_id"},
						{ProtoField: "recipient"},
						{ProtoField: "amount"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
