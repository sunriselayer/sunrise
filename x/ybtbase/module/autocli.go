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
					Use:            "create [admin] [permissioned]",
					Short:          "Send a create tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "admin"}, {ProtoField: "permissioned"}},
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
					RpcMethod:      "GrantYieldPermission",
					Use:            "grant-yield-permission [creator] [target]",
					Short:          "Send a grant-yield-permission tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}, {ProtoField: "target"}},
				},
				{
					RpcMethod:      "RevokeYieldPermission",
					Use:            "revoke-yield-permission [creator] [target]",
					Short:          "Send a revoke-yield-permission tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "creator"}, {ProtoField: "target"}},
				},
				{
					RpcMethod:      "ClaimYield",
					Use:            "claim-yield [token-creator]",
					Short:          "Send a claim-yield tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "token_creator"}},
				},
				{
					RpcMethod:      "UpdateAdmin",
					Use:            "update-admin [new-admin]",
					Short:          "Send a update-admin tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "new_admin"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
