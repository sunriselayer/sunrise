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
				{
					RpcMethod: "AddressBonded",
					Use:       "address-bonded [address]",
					Short:     "Shows the bonded amount of an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "AddressUnbonding",
					Use:       "address-unbonding [address]",
					Short:     "Shows the unbonding of an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "CalculateBondingAmount",
					Use:       "calculate-bonding-amount [validator_address] [share]",
					Short:     "Shows the bonding amount of a validator by share",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "share"},
					},
				},
				{
					RpcMethod: "CalculateShare",
					Use:       "calculate-share [validator_address] [amount]",
					Short:     "Shows the share of a validator by amount",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "ClaimableRewards",
					Use:       "claimable-rewards [address] [validator_address]",
					Short:     "Shows the claimable rewards of a validator",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
						{ProtoField: "validator_address"},
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
					RpcMethod: "NonVotingDelegate",
					Use:       "non-voting-delegate [validator_address] [amount]",
					Short:     "Send a non-voting-delegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "NonVotingUndelegate",
					Use:       "non-voting-undelegate [validator_address] [amount] [recipient?]",
					Short:     "Send a non-voting-undelegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "amount"},
						{ProtoField: "recipient", Optional: true},
					},
				},
				{
					RpcMethod: "ClaimRewards",
					Use:       "claim-rewards [validator_address]",
					Short:     "Send a claim-rewards tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
					},
				},
				{
					RpcMethod: "CreateValidator",
					Use:       "create-validator [validator_address] [min_self_delegation] [amount]",
					Short:     "Send a create-validator tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "min_self_delegation"},
						{ProtoField: "amount"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
