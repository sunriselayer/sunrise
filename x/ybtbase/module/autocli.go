package ybtbase

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/sunriselayer/sunrise/x/ybtbase/types"
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
					RpcMethod:      "Create",
					Use:            "create [admin] [permission-mode]",
					Short:          "Create a new base YBT token",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admin"}, {ProtoField: "permission_mode"}},
				},
				{
					RpcMethod:      "Mint",
					Use:            "mint [creator] [amount]",
					Short:          "Send a mint tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "Burn",
					Use:            "burn [creator] [amount]",
					Short:          "Send a burn tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "AddYield",
					Use:            "add-yield [creator] [amount]",
					Short:          "Send a add-yield tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "GrantPermission",
					Use:            "grant-permission [token-creator] [target]",
					Short:          "Grant permission to an address (whitelist mode only)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "target"}},
				},
				{
					RpcMethod:      "RevokePermission",
					Use:            "revoke-permission [token-creator] [target]",
					Short:          "Revoke permission or block an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "target"}},
				},
				{
					RpcMethod:      "ClaimYield",
					Use:            "claim-yield [token-creator]",
					Short:          "Send a claim-yield tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}},
				},
				{
					RpcMethod:      "UpdateAdmin",
					Use:            "update-admin [token-creator] [new-admin]",
					Short:          "Transfer admin rights for a token",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}, {ProtoField: "new_admin"}},
				},
				{
					RpcMethod:      "Send",
					Use:            "send [to-address] [token-creator] [amount]",
					Short:          "Send base YBT tokens (for restricted tokens)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "to_address"}, {ProtoField: "token_creator"}, {ProtoField: "amount"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
